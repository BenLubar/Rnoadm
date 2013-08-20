goog.provide('rnoadm.hud');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.net');



/**
 * @constructor
 * @struct
 */
rnoadm.hud.HUD = function(name, paint, mouseMoved, click) {
  /**
   * @type {string}
   * @private
   * @const
   */
  this.name_ = name;

  /**
   * @type {function(number, number)}
   * @private
   * @const
   */
  this.paint_ = paint;

  /**
   * @type {function(number, number, number, number):boolean}
   * @private
   * @const
   */
  this.mouseMoved_ = mouseMoved;

  /**
   * @type {function(number, number, number, number):boolean}
   * @private
   * @const
   */
  this.click_ = click;
};


/**
 * @type {Array.<rnoadm.hud.HUD>}
 * @private
 * @const
 */
rnoadm.hud.activeHuds_ = [];


/**
 * @private
 * @struct
 * @const
 */
rnoadm.hud.text_ = {
  charactert: new rnoadm.gfx.Text('Character', '#ccc', true),
  change_name: new rnoadm.gfx.Text('change name', '#aaa', false),
  change_nameh: new rnoadm.gfx.Text('change name', '#fff', false),
  change_skin_color: new rnoadm.gfx.Text('change skin color', '#aaa', false),
  change_skin_colorh: new rnoadm.gfx.Text('change skin color', '#fff', false),
  change_shirt_color: new rnoadm.gfx.Text('change shirt color', '#aaa', false),
  change_shirt_colorh: new rnoadm.gfx.Text('change shirt color', '#fff', false),
  change_pants_color: new rnoadm.gfx.Text('change pants color', '#aaa', false),
  change_pants_colorh: new rnoadm.gfx.Text('change pants color', '#fff', false),
  confirmt: new rnoadm.gfx.Text('Confirm', '#aaa', true),
  confirmth: new rnoadm.gfx.Text('Confirm', '#fff', true),
  push_enter_to_chat: new rnoadm.gfx.Text('[push enter to chat]', '#ccc',
      false, true)
};


/**
 * @type {Object.<string, function(Object):rnoadm.hud.HUD>}
 * @private
 * @const
 */
rnoadm.hud.huds_ = {};

rnoadm.hud.huds_['cc'] = function(data) {
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
    rnoadm.gfx.ctx.fillStyle = 'rgba(0,0,0,.7)';
    rnoadm.gfx.ctx.fillRect(0, 0, w, h);
    w /= rnoadm.gfx.TILE_SIZE * 2;
    h /= rnoadm.gfx.TILE_SIZE * 2;
    sprites.forEach(function(sprite) {
      sprite.paint(w - 3, h + 1);
    });
    rnoadm.hud.text_.charactert.paint(w + 2, h - 4);
    if (hover_name) {
      nameh.paint(w, h + 2);
      rnoadm.hud.text_.change_nameh.paint(w, h + 2.5);
    } else {
      name.paint(w, h + 2);
      rnoadm.hud.text_.change_name.paint(w, h + 2.5);
    }
    if (hover_gender) {
      change_genderh.paint(w + 2, h - 3);
    } else {
      change_gender.paint(w + 2, h - 3);
    }
    if (hover_skin) {
      rnoadm.hud.text_.change_skin_colorh.paint(w + 2, h - 2);
    } else {
      rnoadm.hud.text_.change_skin_color.paint(w + 2, h - 2);
    }
    if (hover_shirt) {
      rnoadm.hud.text_.change_shirt_colorh.paint(w + 2, h - 1);
    } else {
      rnoadm.hud.text_.change_shirt_color.paint(w + 2, h - 1);
    }
    if (hover_pants) {
      rnoadm.hud.text_.change_pants_colorh.paint(w + 2, h);
    } else {
      rnoadm.hud.text_.change_pants_color.paint(w + 2, h);
    }
    if (hover_confirm) {
      rnoadm.hud.text_.confirmth.paint(w, h + 4);
    } else {
      rnoadm.hud.text_.confirmt.paint(w, h + 4);
    }
  }, function(x, y, w, h) {
    x /= rnoadm.gfx.TILE_SIZE;
    y /= rnoadm.gfx.TILE_SIZE;
    w /= rnoadm.gfx.TILE_SIZE * 2;
    h /= rnoadm.gfx.TILE_SIZE * 2;

    if (y >= h + 2 - name.height() && y <= h + 2.5 && Math.abs(w - x) <=
        Math.max(name.width(), rnoadm.hud.text_.change_name.width()) / 2) {
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
    } else if (y >= h - 2 - rnoadm.hud.text_.change_skin_color.height() &&
        y <= h - 2 && Math.abs(w + 2 - x) <=
        rnoadm.hud.text_.change_skin_color.width() / 2) {
      if (!hover_skin)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = true;
      hover_shirt = false;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - 1 - rnoadm.hud.text_.change_shirt_color.height() &&
        y <= h - 1 && Math.abs(w + 2 - x) <=
        rnoadm.hud.text_.change_shirt_color.width() / 2) {
      if (!hover_shirt)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = true;
      hover_pants = false;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h - rnoadm.hud.text_.change_pants_color.height() &&
        y <= h && Math.abs(w + 2 - x) <=
        rnoadm.hud.text_.change_pants_color.width() / 2) {
      if (!hover_pants)
        rnoadm.gfx.repaint();
      hover_name = false;
      hover_gender = false;
      hover_skin = false;
      hover_shirt = false;
      hover_pants = true;
      hover_confirm = false;
      hover_none = false;
    } else if (y >= h + 4 - rnoadm.hud.text_.confirmt.height() &&
        y <= h + 4 && Math.abs(w - x) <=
        rnoadm.hud.text_.confirmt.width() / 2) {
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
  });
};


rnoadm.net.addHandler('HUD', function(data) {
  rnoadm.hud.activeHuds_.length = 0;
  if (data['N'].length) {
    goog.asserts.assert(data['N'] in rnoadm.hud.huds_,
                        'unknown hud: ' + data['N']);
    var hud = rnoadm.hud.huds_[data['N']](data['D'] || {});
    hud.mouseMoved_(rnoadm.hud.lastMouseMove_[0], rnoadm.hud.lastMouseMove_[1],
                    rnoadm.hud.lastMouseMove_[2], rnoadm.hud.lastMouseMove_[3]);
    rnoadm.hud.activeHuds_.push(hud);
  }
  rnoadm.gfx.repaint();
});

rnoadm.gfx.paintHuds = function(w, h) {
  if (rnoadm.hud.chat_ == null) {
    if (rnoadm.hud.activeHuds_.length == 0) {
      rnoadm.hud.text_.push_enter_to_chat.paint(0.1,
          h / rnoadm.gfx.TILE_SIZE - 0.1);
    }
  } else {
    new rnoadm.gfx.Text(rnoadm.hud.chat_ + '_', '#fff', false, true).paint(
        0.1, h / rnoadm.gfx.TILE_SIZE - 0.1);
  }
  rnoadm.hud.messages_.forEach(function(message, i) {
    message.paint(0.1, h / rnoadm.gfx.TILE_SIZE - i / 2 - 0.6);
  });
  rnoadm.hud.activeHuds_.forEach(function(hud) {
    hud.paint_(w, h);
  });
};


/**
 * @type {Array.<number>}
 * @private
 */
rnoadm.hud.lastMouseMove_ = [-Infinity, -Infinity, 1, 1];

rnoadm.gfx.mouseMovedHud = function(x, y, w, h) {
  rnoadm.hud.lastMouseMove_ = [x, y, w, h];
  for (var i = rnoadm.hud.activeHuds_.length - 1; i >= 0; i--) {
    if (rnoadm.hud.activeHuds_[i].mouseMoved_(x, y, w, h))
      return true;
  }
  return false;
};

rnoadm.gfx.clickHud = function(x, y, w, h) {
  rnoadm.gfx.mouseMovedHud(x, y, w, h);
  for (var i = rnoadm.hud.activeHuds_.length - 1; i >= 0; i--) {
    if (rnoadm.hud.activeHuds_[i].click_(x, y, w, h))
      return true;
  }
  return false;
};

rnoadm.net.addHandler('Inventory', function(inventory) {
  window['console']['log'](inventory);
});


window.addEventListener('keydown', function(e) {
  if (e.ctrlKey || e.altKey) return;
  if (e.keyCode < 20) e.preventDefault();
  switch (e.keyCode) {
  case 8: // backspace
    if (rnoadm.hud.chat_ !== null && rnoadm.hud.chat_.length > 0) {
      rnoadm.hud.chat_ = rnoadm.hud.chat_.substring(0,
        rnoadm.hud.chat_.length - 1);
      rnoadm.gfx.repaint();
    }
    break;
  case 13: // enter
    if (rnoadm.hud.chat_ === null) {
      rnoadm.hud.chat_ = '';
    } else {
      rnoadm.net.send({'Chat': rnoadm.hud.chat_});
      rnoadm.hud.chat_ = null;
    }
    rnoadm.gfx.repaint();
    break;
  case 27: // esc
    if (rnoadm.hud.chat_ !== null) {
      rnoadm.hud.chat_ = null;
      rnoadm.gfx.repaint();
    }
    break;
  }
}, false);


window.addEventListener('keypress', function(e) {
  if (rnoadm.hud.chat_ !== null) {
    rnoadm.hud.chat_ += String.fromCharCode(e.charCode);
    rnoadm.gfx.repaint();
  }
}, false);


rnoadm.net.onConnect.push(function() {
  rnoadm.hud.chat_ = null;
});


/**
 * @type {?string}
 * @private
 */
rnoadm.hud.chat_ = null;


/**
 * @type {Array.<rnoadm.gfx.Text>}
 * @private
 * @const
 */
rnoadm.hud.messages_ = [];

rnoadm.net.addHandler('Msg', function(messages) {
  messages.forEach(function(message) {
    rnoadm.hud.messages_.unshift(new rnoadm.gfx.Text(message['T'], message['C'], false, true));
    window.setTimeout(function() {
      rnoadm.hud.messages_.pop();
    }, 60000);
    rnoadm.gfx.repaint();
  });
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
