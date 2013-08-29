goog.provide('rnoadm.hud.forge');

goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');

rnoadm.hud.forge = function(data) {
  var forgeSprite = new rnoadm.gfx.Sprite('forge', '#888', '', 0, 0, 64, 32);
  var forgeText = new rnoadm.gfx.Text('Forge', '#ccc', true);

  return new rnoadm.hud.HUD('forge', function(w, h) {
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(w/2 - rnoadm.gfx.TILE_SIZE * 7.5,
                            h/2 - rnoadm.gfx.TILE_SIZE * 6,
                            rnoadm.gfx.TILE_SIZE * 16,
                            rnoadm.gfx.TILE_SIZE * 12);
    rnoadm.gfx.ctx.fillStyle = '#000';
    rnoadm.gfx.ctx.fillRect(w/2 - rnoadm.gfx.TILE_SIZE * 7.5,
                            h/2 - rnoadm.gfx.TILE_SIZE * 6,
                            rnoadm.gfx.TILE_SIZE * 16,
                            rnoadm.gfx.TILE_SIZE);
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;

    forgeSprite.paint(w2 - 7, h2 - 5);
    forgeText.paint(w2 + 1, h2 - 5);
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2) > 7.5 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      return false;
    }
    // TODO: hover
    return true;
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2) > 7.5 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      rnoadm.hud.hide('forge');
      return false;
    }
    // TODO: click
    return true;
  }, function(x, y, w, h) {
    var w2 = w / rnoadm.gfx.TILE_SIZE / 2;
    var h2 = h / rnoadm.gfx.TILE_SIZE / 2;
    if (Math.abs(x / rnoadm.gfx.TILE_SIZE - w2) > 7.5 ||
        Math.abs(y / rnoadm.gfx.TILE_SIZE - h2) > 6) {
      rnoadm.hud.hide('forge');
      return false;
    }
    // TODO: right click
    return true;
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
