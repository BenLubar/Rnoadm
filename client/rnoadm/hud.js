goog.provide('rnoadm.hud');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.net');
goog.require('rnoadm.state.Object');



/**
 * @constructor
 * @struct
 */
rnoadm.hud.HUD = function(name, paint, mouseMoved, click, rightClick) {
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

  /**
   * @type {function(number, number, number, number):boolean}
   * @private
   * @const
   */
  this.rightClick_ = rightClick;
};


/**
 * @type {Array.<rnoadm.hud.HUD>}
 * @private
 * @const
 */
rnoadm.hud.activeHuds_ = [];


/**
 * @struct
 * @const
 */
rnoadm.hud.text = {
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


/**
 * @param {string} id
 * @param {function(Object):rnoadm.hud.HUD} f
 */
rnoadm.hud.register = function(id, f) {
  goog.asserts.assert(!(id in rnoadm.hud.huds_), 'double hud register: ' + id);
  rnoadm.hud.huds_[id] = f;
};


/**
 * @param {string} id
 * @param {Object=} opt_data
 */
rnoadm.hud.show = function(id, opt_data) {
  goog.asserts.assert(id in rnoadm.hud.huds_,
                      'unknown hud: ' + id);
  var hud = rnoadm.hud.huds_[id](opt_data || {});
  hud.mouseMoved_(rnoadm.hud.lastMouseMove_[0], rnoadm.hud.lastMouseMove_[1],
                  rnoadm.hud.lastMouseMove_[2], rnoadm.hud.lastMouseMove_[3]);
  var found = false;
  rnoadm.hud.activeHuds_.forEach(function(h, i) {
    if (h.name_ == id) {
      found = true;
      rnoadm.hud.activeHuds_[i] = hud;
      return false;
    }
  });
  if (!found) {
    rnoadm.hud.activeHuds_.push(hud);
  }
  rnoadm.gfx.repaint();
};


/**
 * @param {string} id
 */
rnoadm.hud.hide = function(id) {
  rnoadm.hud.activeHuds_.forEach(function(h, i) {
    if (h.name_ == id) {
      rnoadm.hud.activeHuds_.splice(i, 1);
      rnoadm.gfx.repaint();
      return false;
    }
  });
};

rnoadm.net.addHandler('HUD', function(data) {
  rnoadm.hud.activeHuds_.length = 0;
  rnoadm.gfx.repaint();
  if (data['N'].length) {
    rnoadm.hud.show(data['N'], data['D']);
  }
});


/**
 * @type {boolean}
 * @private
 */
rnoadm.hud.hover_inventory_ = false;


/**
 * @type {rnoadm.gfx.Sprite}
 * @private
 * @const
 */
rnoadm.hud.inventoryIcon_ = new rnoadm.gfx.Sprite('ui_icons', '#888', '',
                                                  0, 0, 32, 32);


/**
 * @type {rnoadm.gfx.Sprite}
 * @private
 * @const
 */
rnoadm.hud.inventoryIconh_ = new rnoadm.gfx.Sprite('ui_icons', '#bbb', '',
                                                   0, 0, 32, 32);


rnoadm.gfx.paintHuds = function(w, h) {
  if (rnoadm.hud.chat_ == null) {
    rnoadm.hud.text.push_enter_to_chat.paint(0.1,
        h / rnoadm.gfx.TILE_SIZE - 0.1);
  } else {
    new rnoadm.gfx.Text(rnoadm.hud.chat_ + '_', '#fff', false, true).paint(
        0.1, h / rnoadm.gfx.TILE_SIZE - 0.1);
  }
  rnoadm.hud.messages_.forEach(function(message, i) {
    message.paint(0.1, h / rnoadm.gfx.TILE_SIZE - i / 2 - 0.6);
  });
  if (rnoadm.hud.hover_inventory_) {
    rnoadm.hud.inventoryIconh_.paint(w / rnoadm.gfx.TILE_SIZE - 1,
                                     h / rnoadm.gfx.TILE_SIZE);
  } else {
    rnoadm.hud.inventoryIcon_.paint(w / rnoadm.gfx.TILE_SIZE - 1,
                                    h / rnoadm.gfx.TILE_SIZE);
  }
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
  if (rnoadm.hud.activeHuds_.length == 0 && x >= w - rnoadm.gfx.TILE_SIZE &&
      x <= w && y >= h - rnoadm.gfx.TILE_SIZE && y <= h) {
    if (!rnoadm.hud.hover_inventory_) {
      rnoadm.hud.hover_inventory_ = true;
      rnoadm.gfx.repaint();
    }
    return true;
  } else if (rnoadm.hud.hover_inventory_) {
    rnoadm.hud.hover_inventory_ = false;
    rnoadm.gfx.repaint();
  }
  for (var i = rnoadm.hud.activeHuds_.length - 1; i >= 0; i--) {
    if (rnoadm.hud.activeHuds_[i].mouseMoved_(x, y, w, h))
      return true;
  }
  return false;
};

rnoadm.gfx.clickHud = function(x, y, w, h) {
  rnoadm.gfx.mouseMovedHud(x, y, w, h);
  if (rnoadm.hud.hover_inventory_) {
    rnoadm.hud.show('inv');
    return true;
  }
  for (var i = rnoadm.hud.activeHuds_.length - 1; i >= 0; i--) {
    if (rnoadm.hud.activeHuds_[i].click_(x, y, w, h))
      return true;
  }
  return false;
};


rnoadm.gfx.rightClickHud = function(x, y, w, h) {
  rnoadm.gfx.mouseMovedHud(x, y, w, h);
  for (var i = rnoadm.hud.activeHuds_.length - 1; i >= 0; i--) {
    if (rnoadm.hud.activeHuds_[i].rightClick_(x, y, w, h))
      return true;
  }
  return false;
};


window.addEventListener('keydown', function(e) {
  if (!rnoadm.gfx.canvas.parentNode) return;
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
