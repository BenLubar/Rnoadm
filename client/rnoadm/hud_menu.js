goog.provide('rnoadm.hud.menu');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');
goog.require('rnoadm.net');
goog.require('rnoadm.state.Object');


rnoadm.hud.register('menu', function(data) {
  /** @type {rnoadm.state.Object} */
  var object = goog.asserts.assertInstanceof(data, rnoadm.state.Object);

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

  /** @type {Array.<string>} */
  var actions = object.actions;

  actions.forEach(function(action) {
    action = action + ' ' + object.name;
    options.push(new rnoadm.gfx.Text(action, '#aaa', false, true));
    optionsh.push(new rnoadm.gfx.Text(action, '#fff', false, true));
  });

  /** @type {number} */
  var width = -2;
  options.forEach(function(option) {
    width = Math.max(option.width(), width);
  });
  width = (width + 0.2) * rnoadm.gfx.TILE_SIZE;
  /** @type {number} */
  var height = Math.floor(rnoadm.gfx.TILE_SIZE * (actions.length / 2 +
      0.2));

  return new rnoadm.hud.HUD('menu', function(w, h) {
    if (isInitial) {
      isInitial = false;
      if (initialX + width > w)
        initialX = w - width;
      if (initialY + height > h)
        initialY = h - height;
    }
    rnoadm.gfx.ctx.fillStyle = 'rgba(0,0,0,.7)';
    rnoadm.gfx.ctx.fillRect(initialX || 0, initialY || 0, width, height);
    if (hover != -1) {
      rnoadm.gfx.ctx.fillStyle = '#000';
      rnoadm.gfx.ctx.fillRect(initialX || 0, (initialY || 0) + (hover / 2 +
          0.1) * rnoadm.gfx.TILE_SIZE, width, rnoadm.gfx.TILE_SIZE / 2);
    }
    var x = initialX / rnoadm.gfx.TILE_SIZE + 0.1;
    var y = initialY / rnoadm.gfx.TILE_SIZE + 0.5;
    for (var i = 0; i < actions.length; i++) {
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
        0.1 * rnoadm.gfx.TILE_SIZE) && y < Math.floor(initialY +
        (actions.length / 2 + 0.1) * rnoadm.gfx.TILE_SIZE)) {
      hover = Math.floor((y - initialY) * 2 / rnoadm.gfx.TILE_SIZE - 0.1);
      if (hover >= actions.length)
        hover = -1;
    } else {
      hover = -1;
      if (x < initialX - rnoadm.gfx.TILE_SIZE || x > initialX + width +
          rnoadm.gfx.TILE_SIZE || y < Math.floor(initialY -
          rnoadm.gfx.TILE_SIZE) || y > initialY + height +
          rnoadm.gfx.TILE_SIZE) {
        rnoadm.hud.hide('menu');
      }
    }
    rnoadm.gfx.repaint();
  }, function(x, y, w, h) {
    if (hover != -1)
      rnoadm.net.send({'Interact': {
        'ID': object.id,
        'X': object.x,
        'Y': object.y,
        'Action': actions[hover]
      }});
    rnoadm.hud.hide('menu');
    return true;
  }, function(x, y, w, h) {
    rnoadm.hud.hide('menu');
    return true;
  });
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
