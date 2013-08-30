goog.provide('rnoadm.hud.inv');
goog.provide('rnoadm.state.inventory');

goog.require('rnoadm.gfx');
goog.require('rnoadm.hud');
goog.require('rnoadm.net');
goog.require('rnoadm.state.Object');


rnoadm.hud.inv = function(data) {
  /** @type {number} */
  var hoverX = -1;
  /** @type {number} */
  var hoverY = -1;
  /** @type {number} */
  var mouseX = -Infinity;
  /** @type {number} */
  var mouseY = -Infinity;

  return new rnoadm.hud.HUD('inv', function(w, h) {
    var rows = Math.ceil(rnoadm.state.inventory.length / 8);
    var xOffset = w / rnoadm.gfx.TILE_SIZE - 8.1;
    var yOffset = h / rnoadm.gfx.TILE_SIZE - 0.1 - rows;
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(w - 8.2 * rnoadm.gfx.TILE_SIZE,
                            h - (rows + 1.2) * rnoadm.gfx.TILE_SIZE,
                            8.2 * rnoadm.gfx.TILE_SIZE,
                            (rows + 0.2) * rnoadm.gfx.TILE_SIZE);
    rnoadm.state.inventory.forEach(function(item) {
      item.paint(xOffset, yOffset);
    });
    if (hoverX >= 0 && hoverX < 8 && hoverY >= 0 && hoverY < rows &&
        hoverX + hoverY * 8 < rnoadm.state.inventory.length) {
      var text = new rnoadm.gfx.Text(rnoadm.state.inventory[hoverX +
          hoverY * 8].name, '#fff', false, true);
      var width = (text.width() + 0.2) * rnoadm.gfx.TILE_SIZE;
      var height = (text.height() + 0.2) * rnoadm.gfx.TILE_SIZE;
      rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
      rnoadm.gfx.ctx.fillRect(mouseX - width, mouseY, width, height);
      text.paint(mouseX / rnoadm.gfx.TILE_SIZE - text.width() - 0.1,
                 mouseY / rnoadm.gfx.TILE_SIZE + text.height() + 0.1);
    }
  }, function(x, y, w, h) {
    var rows = Math.ceil(rnoadm.state.inventory.length / 8);
    mouseX = x;
    mouseY = y;
    hoverX = Math.floor((x - w) / rnoadm.gfx.TILE_SIZE + 8.1);
    hoverY = Math.floor((y - h) / rnoadm.gfx.TILE_SIZE + rows + 1.1);
    rnoadm.gfx.repaint();
    if (hoverY > rows || hoverX < -1) {
      rnoadm.hud.hide('inv');
    } else {
      return true;
    }
  }, function(x, y, w, h) {
    return true;
  }, function(x, y, w, h) {
    var rows = Math.ceil(rnoadm.state.inventory.length / 8);
    if (hoverX >= 0 && hoverX < 8 && hoverY >= 0 && hoverY < rows &&
        hoverX + hoverY * 8 < rnoadm.state.inventory.length) {
      rnoadm.hud.show('menu', rnoadm.state.inventory[hoverX + hoverY * 8]);
    }
  });
};


/**
 * @type {Array.<rnoadm.state.Object>}
 * @const
 */
rnoadm.state.inventory = [];

rnoadm.net.addHandler('Inventory', function(inventory) {
  /** @type {Object.<string, rnoadm.state.Object>} */
  var old = {};
  rnoadm.state.inventory.forEach(function(item) {
    old[item.id] = item;
  });

  rnoadm.state.inventory.length = 0;
  inventory.forEach(function(item, i) {
    /** @type {rnoadm.state.Object} */
    var obj = old[item['I']];
    if (obj) {
      obj.update(item['O']);
      obj.move(i % 8, Math.floor(i / 8));
    } else {
      obj = new rnoadm.state.Object(i % 8, Math.floor(i / 8),
                                    item['I'], item['O']);
    }
    rnoadm.state.inventory.push(obj);
  });

  rnoadm.gfx.repaint();
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
