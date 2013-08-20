goog.provide('rnoadm.hud.inv');

goog.require('rnoadm.hud');
goog.require('rnoadm.gfx');
goog.require('rnoadm.state.Object');
goog.require('rnoadm.net');


rnoadm.hud.register('inv', function(data) {
  return new rnoadm.hud.HUD('inv', function(w, h) {
    var rows = Math.ceil(rnoadm.hud.inventory_.length / 8);
    var xOffset = w / rnoadm.gfx.TILE_SIZE - 8.1;
    var yOffset = h / rnoadm.gfx.TILE_SIZE - 1.1 - rows;
    rnoadm.gfx.ctx.fillStyle = 'rgba(0,0,0,.7)';
    rnoadm.gfx.ctx.fillRect(w - 8.2 * rnoadm.gfx.TILE_SIZE,
                            h - (rows + 1.2) * rnoadm.gfx.TILE_SIZE,
                            8.2, rows + 0.2);
    rnoadm.hud.inventory_.forEach(function(item) {
      item.paint(xOffset, yOffset);
    });
  }, function(x, y, w, h) {
    // TODO: mouse moved
  }, function(x, y, w, h) {
    // TODO: clicked
  });
});

/**
 * @type {Array.<rnoadm.state.Object>}
 * @private
 * @const
 */
rnoadm.hud.inventory_ = [];

rnoadm.net.addHandler('Inventory', function(inventory) {
  /** @type {Object.<string, rnoadm.state.Object>} */
  var old = {};
  rnoadm.hud.inventory_.forEach(function(item) {
    old[item.id] = item;
  });

  rnoadm.hud.inventory_.length = 0;
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
    rnoadm.hud.inventory_.push(obj);
  });

  rnoadm.gfx.repaint();
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
