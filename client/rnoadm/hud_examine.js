goog.provide('rnoadm.hud.examine');

goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');

rnoadm.hud.register('examine', function(data) {
  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var name = new rnoadm.gfx.Text(data['N'], '#ccc', true, true);

  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var examine = new rnoadm.gfx.Text(data['E'], '#ccc', false, true);

  /**
   * @type {Array.<Array.<rnoadm.gfx.Text>>}
   * @const
   */
  var info = [];

  (data['I'] || []).forEach(function(line) {
    /** @type {Array.<rnoadm.gfx.Text>} */
    var infoLine = [];
    info.push(infoLine);
    line.forEach(function(text) {
      infoLine.push(new rnoadm.gfx.Text(text[0], text[1], false, true));
    });
  });

  return new rnoadm.hud.HUD('examine', function(w, h) {
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(0, 0, w, h);
    var x = w / 2 / rnoadm.gfx.TILE_SIZE - 6;
    var y = h / 2 / rnoadm.gfx.TILE_SIZE - 6;
    name.paint(x, y);
    y += 0.5;
    examine.paint(x, y);
    y += 1;
    info.forEach(function(line) {
      x = w / 2 / rnoadm.gfx.TILE_SIZE - 6;
      line.forEach(function(text) {
        text.paint(x, y);
        x += text.width();
      });
      y += 0.5;
    });
  }, function(x, y, w, h) {
    return true;
  }, function(x, y, w, h) {
    rnoadm.hud.hide('examine');
    return true;
  }, function(x, y, w, h) {
    rnoadm.hud.hide('examine');
    return true;
  });
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
