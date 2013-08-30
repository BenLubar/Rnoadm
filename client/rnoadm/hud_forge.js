goog.provide('rnoadm.hud.forge');

goog.require('rnoadm.state.inventory');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');

rnoadm.hud.forge = function(data) {
  var ores = [];
  data['O'].forEach(function(id) {
    rnoadm.state.inventory.forEach(function(item) {
      if (item.id == id) {
        ores.push(item.clone());
      }
    });
  });
  ores.forEach(function(ore, i) {
    ore.x = ore.prevX = i % 8;
    ore.y = ore.prevY = Math.floor(i / 8);
  });

  var forgeSprite = new rnoadm.gfx.Sprite('forge', 'no', '', 0, 0, 64, 32);
  var forgeText = new rnoadm.gfx.Text('Forge', '#ccc', true, true);

  return new rnoadm.hud.HUD('forge', function(w, h) {
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(w/2 - rnoadm.gfx.TILE_SIZE * 7.5,
                            h/2 - rnoadm.gfx.TILE_SIZE * 5,
                            rnoadm.gfx.TILE_SIZE * 16,
                            rnoadm.gfx.TILE_SIZE * 11);
    rnoadm.gfx.ctx.fillStyle = '#000';
    rnoadm.gfx.ctx.beginPath();
    rnoadm.gfx.ctx.moveTo(w/2 - rnoadm.gfx.TILE_SIZE * 7.5,
      h/2 - rnoadm.gfx.TILE_SIZE * 5);
    rnoadm.gfx.ctx.lineTo(w/2 + rnoadm.gfx.TILE_SIZE * 8.5,
      h/2 - rnoadm.gfx.TILE_SIZE * 5);
    rnoadm.gfx.ctx.lineTo(w/2 + rnoadm.gfx.TILE_SIZE * 7,
      h/2 - rnoadm.gfx.TILE_SIZE * 6);
    rnoadm.gfx.ctx.lineTo(w/2 - rnoadm.gfx.TILE_SIZE * 7.5,
      h/2 - rnoadm.gfx.TILE_SIZE * 6);
    rnoadm.gfx.ctx.closePath();
    rnoadm.gfx.ctx.fill();

    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;

    forgeSprite.paint(w2 - 7, h2 - 5);
    forgeText.paint(w2 - 5, h2 - 5.1);

    ores.forEach(function(ore) {
      ore.paint(w2 + 0.5, h2 - 4);
    });
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2 - 0.5) > 8 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      return false;
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
