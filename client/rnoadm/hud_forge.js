goog.provide('rnoadm.hud.forge');

goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');
goog.require('rnoadm.net');
goog.require('rnoadm.state.inventory');

rnoadm.hud.forge = function(data) {
  var ores = [];
  data['O'].forEach(function(itemData) {
    rnoadm.state.inventory.forEach(function(item) {
      if (item.id == itemData['i']) {
        ores.push(item.clone());
      }
    });
  });
  ores.forEach(function(ore, i) {
    ore.x = ore.prevX = i % 8;
    ore.y = ore.prevY = Math.floor(i / 8);
  });

  var contents = [];
  data['C'].forEach(function(sprites) {
    contents.push(sprites.map(rnoadm.gfx.Sprite.fromNetwork));
  });

  var done = new Date(data['T']);

  var hoverItem = null;

  var forgeSprite = new rnoadm.gfx.Sprite('forge', 'no', '', 0, 0, 64, 32);
  var forgeText = new rnoadm.gfx.Text('Forge', '#ccc', true, true);
  var inventoryText = new rnoadm.gfx.Text('Inventory', '#aaa', false, true);
  var itemHoverText = null;
  var itemHoverText2 = null;

  return new rnoadm.hud.HUD('forge', function(w, h) {
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(w / 2 - rnoadm.gfx.TILE_SIZE * 7.5,
                            h / 2 - rnoadm.gfx.TILE_SIZE * 5,
                            rnoadm.gfx.TILE_SIZE * 16,
                            rnoadm.gfx.TILE_SIZE * 11);
    rnoadm.gfx.ctx.fillStyle = '#000';
    rnoadm.gfx.ctx.beginPath();
    rnoadm.gfx.ctx.moveTo(w / 2 - rnoadm.gfx.TILE_SIZE * 7.5,
        h / 2 - rnoadm.gfx.TILE_SIZE * 5);
    rnoadm.gfx.ctx.lineTo(w / 2 + rnoadm.gfx.TILE_SIZE * 8.5,
        h / 2 - rnoadm.gfx.TILE_SIZE * 5);
    rnoadm.gfx.ctx.lineTo(w / 2 + rnoadm.gfx.TILE_SIZE * 7,
        h / 2 - rnoadm.gfx.TILE_SIZE * 6);
    rnoadm.gfx.ctx.lineTo(w / 2 - rnoadm.gfx.TILE_SIZE * 7.25,
        h / 2 - rnoadm.gfx.TILE_SIZE * 6);
    rnoadm.gfx.ctx.closePath();
    rnoadm.gfx.ctx.fill();

    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;

    forgeSprite.paint(w2 - 7, h2 - 5);
    forgeText.paint(w2 - 5.5, h2 - 5.1);

    var remaining = done - new Date();
    if (remaining > 0) {
      rnoadm.gfx.repaint(remaining % 1000);
      new rnoadm.gfx.Text(Math.floor(remaining / 1000) + 's remaining',
          '#fff', false, true).paint(w2 - 7.5, h2 - 4.6);
    }

    contents.forEach(function(sprites, i) {
      sprites.forEach(function(sprite) {
        sprite.paint(w2 - 7.5 + i % 8, h2 - 3.5 + Math.floor(i / 8));
      });
    });

    inventoryText.paint(w2 + 0.5, h2 - 4.6);
    ores.forEach(function(ore) {
      ore.paint(w2 + 0.5, h2 - 3.5);
    });

    if (itemHoverText) {
      itemHoverText.paint(w2 + 0.5, h2 + 5);
    }
    if (itemHoverText2) {
      itemHoverText2.paint(w2 + 0.5, h2 + 5.5);
    }
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    x = x / rnoadm.gfx.TILE_SIZE - w2;
    y = y / rnoadm.gfx.TILE_SIZE - h2;
    if (Math.abs(x - 0.5) > 8 ||
        Math.abs(y) > 6) {
      return false;
    }
    var inventoryX = Math.floor(x - 0.5);
    var inventoryY = Math.floor(y + 4.5);
    if (inventoryX >= 0 && inventoryX < 8 && inventoryY >= 0 &&
        inventoryX + inventoryY * 8 < ores.length) {
      var ore = ores[inventoryX + inventoryY * 8];
      if (ore != hoverItem) {
        var dat = data['O'][inventoryX + inventoryY * 8];
        hoverItem = ore;
        itemHoverText = new rnoadm.gfx.Text(ore.name, '#ccc', false, true);
        itemHoverText2 = new rnoadm.gfx.Text('quality ' + dat['q'] +
            '  weight ' + (Math.round(dat['w'] / 100) / 10) + 'kg' +
            '  volume ' + dat['v'] + 'cc', '#aaa', false, true);
        rnoadm.gfx.repaint();
      }
    } else {
      if (hoverItem) {
        hoverItem = null;
        itemHoverText = null;
        itemHoverText2 = null;
        rnoadm.gfx.repaint();
      }
    }
    // TODO: hover
    return true;
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2 - 0.5) > 8 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      rnoadm.hud.hide('forge');
      return false;
    }
    if (hoverItem) {
      rnoadm.net.send({'HUD': {'N': 'forge', 'D': {'A': 'a', 'I': hoverItem.id}}});
      return true;
    }
    // TODO: click
    return true;
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2 - 0.5) > 8 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      rnoadm.hud.hide('forge');
      return false;
    }
    // TODO: right click
    return true;
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
