goog.provide('rnoadm.hud.menu2');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');
goog.require('rnoadm.net');
goog.require('rnoadm.state.Object');


rnoadm.hud.menu2 = function(data) {
  /** @type {number} */
  var hover = -1;
  /** @type {?number} */
  var initialX = null;
  /** @type {?number} */
  var initialY = null;
  /** @type {boolean} */
  var isInitial = false;

  /** @type {Array.<rnoadm.gfx.Text>} */
  var options = [];
  /** @type {Array.<rnoadm.gfx.Text>} */
  var optionsh = [];

  goog.asserts.assertArray(data).forEach(function(object) {
    var name = goog.asserts.assertInstanceof(object, rnoadm.state.Object).name;
    options.push(new rnoadm.gfx.Text(name, '#aaa', false, true));
    optionsh.push(new rnoadm.gfx.Text(name, '#fff', false, true));
  });

  /** @type {number} */
  var width = -2;
  options.forEach(function(option) {
    width = Math.max(option.width(), width);
  });
  width = (width + 0.2) * rnoadm.gfx.TILE_SIZE;
  /** @type {number} */
  var height = Math.floor(rnoadm.gfx.TILE_SIZE * (data.length / 2 + 0.2));

  return new rnoadm.hud.HUD('menu2', function(w, h) {
    if (isInitial) {
      isInitial = false;
      if (initialX + width > w)
        initialX = w - width;
      if (initialY + height > h)
        initialY = h - height;
    }
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(initialX || 0, initialY || 0, width, height);
    if (hover != -1) {
      rnoadm.gfx.ctx.fillStyle = '#000';
      rnoadm.gfx.ctx.fillRect(initialX || 0, (initialY || 0) + (hover / 2 +
          0.1) * rnoadm.gfx.TILE_SIZE, width, rnoadm.gfx.TILE_SIZE / 2);
    }
    var x = initialX / rnoadm.gfx.TILE_SIZE + 0.1;
    var y = initialY / rnoadm.gfx.TILE_SIZE + 0.5;
    for (var i = 0; i < data.length; i++) {
      if (i == hover) {
        optionsh[i].paint(x, y + i / 2);
      } else {
        options[i].paint(x, y + i / 2);
      }
    }
  }, function(x, y, w, h) {
    if (initialX === null) {
      initialX = x;
      initialY = y;
      isInitial = true;
      rnoadm.gfx.repaint();
    }
    if (x >= initialX && x < initialX + width && y >= Math.floor(initialY +
        0.1 * rnoadm.gfx.TILE_SIZE) && y < initialY + Math.floor((data.length /
        2 + 0.1) * rnoadm.gfx.TILE_SIZE)) {
      hover = Math.floor((y - initialY) * 2 / rnoadm.gfx.TILE_SIZE - 0.1);
      if (hover >= data.length)
        hover = -1;
    } else {
      hover = -1;
      if (x < initialX - rnoadm.gfx.TILE_SIZE || x > initialX + width +
          rnoadm.gfx.TILE_SIZE || y < initialY - rnoadm.gfx.TILE_SIZE ||
          y > initialY + height + rnoadm.gfx.TILE_SIZE) {
        rnoadm.hud.hide('menu2');
      }
    }
    rnoadm.gfx.repaint();
    return true;
  }, function(x, y, w, h) {
    if (hover != -1) {
      rnoadm.hud.show('menu', data[hover]);
    }
    rnoadm.hud.hide('menu2');
    return true;
  }, function(x, y, w, h) {
    rnoadm.hud.hide('menu2');
    return true;
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
