goog.provide('rnoadm.hud.cc');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.hud');
goog.require('rnoadm.net');



rnoadm.hud.cc = function(data) {
  /**
   * @type {Array.<rnoadm.gfx.Sprite>}
   * @const
   */
  var sprites = [];

  data['S'].forEach(function(sprite) {
    sprites.push(rnoadm.gfx.Sprite.fromNetwork(sprite));
  });

  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var name = new rnoadm.gfx.Text(data['N'], '#aaa', true);

  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var nameh = new rnoadm.gfx.Text(data['N'], '#fff', true);

  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var change_gender = new rnoadm.gfx.Text('change gender (' + data['G'] + ')',
                                          '#aaa', false);

  /**
   * @type {rnoadm.gfx.Text}
   * @const
   */
  var change_genderh = new rnoadm.gfx.Text('change gender (' + data['G'] + ')',
                                           '#fff', false);

  /** @type {boolean} */
  var hover_name = false;

  /** @type {boolean} */
  var hover_gender = false;

  /** @type {boolean} */
  var hover_skin = false;

  /** @type {boolean} */
  var hover_shirt = false;

  /** @type {boolean} */
  var hover_pants = false;

  /** @type {boolean} */
  var hover_confirm = false;

  /** @type {boolean} */
  var hover_none = true;

  return new rnoadm.hud.HUD('cc', function(w, h) {
    rnoadm.gfx.ctx.fillStyle = rnoadm.gfx.INTERFACE_FILL;
    rnoadm.gfx.ctx.fillRect(0, 0, w, h);
    w /= rnoadm.gfx.TILE_SIZE * 2;
    h /= rnoadm.gfx.TILE_SIZE * 2;
    sprites.forEach(function(sprite) {
      sprite.paint(w - 3, h + 1);
    });
    rnoadm.hud.text.charactert.paint(w + 2, h - 4);
    if (hover_name) {
      nameh.paint(w, h + 2);
      rnoadm.hud.text.change_nameh.paint(w, h + 2.5);
    } else {
      name.paint(w, h + 2);
      rnoadm.hud.text.change_name.paint(w, h + 2.5);
    }
    if (hover_gender) {
      change_genderh.paint(w + 2, h - 3);
    } else {
      change_gender.paint(w + 2, h - 3);
    }
    if (hover_skin) {
      rnoadm.hud.text.change_skin_colorh.paint(w + 2, h - 2);
    } else {
      rnoadm.hud.text.change_skin_color.paint(w + 2, h - 2);
    }
    if (hover_shirt) {
      rnoadm.hud.text.change_shirt_colorh.paint(w + 2, h - 1);
    } else {
      rnoadm.hud.text.change_shirt_color.paint(w + 2, h - 1);
    }
    if (hover_pants) {
      rnoadm.hud.text.change_pants_colorh.paint(w + 2, h);
    } else {
      rnoadm.hud.text.change_pants_color.paint(w + 2, h);
    }
    if (hover_confirm) {
      rnoadm.hud.text.confirmth.paint(w, h + 4);
    } else {
      rnoadm.hud.text.confirmt.paint(w, h + 4);
    }
  }, function(x, y, w, h) {
    x /= rnoadm.gfx.TILE_SIZE;
    y /= rnoadm.gfx.TILE_SIZE;
    w /= rnoadm.gfx.TILE_SIZE * 2;
    h /= rnoadm.gfx.TILE_SIZE * 2;

    if (y >= h + 2 - name.height() && y <= h + 2.5 && Math.abs(w - x) <=
        Math.max(name.width(), rnoadm.hud.text.change_name.width()) / 2) {
      if (!hover_name)
        rnoadm.gfx.repaint();
      hover_name = true;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - 3 - change_gender.height() && y <= h - 3 &&
        Math.abs(w + 2 - x) <= change_gender.width() / 2) {
      if (!hover_gender)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = true;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - 2 - rnoadm.hud.text.change_skin_color.height() &&
        y <= h - 2 && Math.abs(w + 2 - x) <=
        rnoadm.hud.text.change_skin_color.width() / 2) {
      if (!hover_skin)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = true;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - 1 - rnoadm.hud.text.change_shirt_color.height() &&
        y <= h - 1 && Math.abs(w + 2 - x) <=
        rnoadm.hud.text.change_shirt_color.width() / 2) {
      if (!hover_shirt)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = true;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - rnoadm.hud.text.change_pants_color.height() &&
        y <= h && Math.abs(w + 2 - x) <=
        rnoadm.hud.text.change_pants_color.width() / 2) {
      if (!hover_pants)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = true;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h + 4 - rnoadm.hud.text.confirmt.height() &&
        y <= h + 4 && Math.abs(w - x) <=
        rnoadm.hud.text.confirmt.width() / 2) {
      if (!hover_confirm)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = true;
      hover_none = false;
    } else {
      if (!hover_none)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = false;
      hover_none = true;
    }

    // No interaction with anything other than this hud while it's up.
    return true;
  }, function(x, y, w, h) {
    switch (true) {
      case hover_name:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'name'}});
        break;
      case hover_gender:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'gender'}});
        break;
      case hover_skin:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'skin'}});
        break;
      case hover_shirt:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'shirt'}});
        break;
      case hover_pants:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'pants'}});
        break;
      case hover_confirm:
        rnoadm.net.send({'HUD': {'N': 'cc', 'D': 'accept'}});
        break;
      case hover_none:
        break;
    }

    // No interaction with anything other than this hud while it's up.
    return true;
  }, function(x, y, w, h) {
    // No interaction with anything other than this hud while it's up.
    return true;
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
