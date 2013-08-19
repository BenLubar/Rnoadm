goog.provide('rnoadm.hud');

goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.net');


/**
 * @constructor
 * @struct
 */
rnoadm.hud.HUD = function(name, paint) {
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
  change_name: new rnoadm.gfx.Text('change name', '#aaa', false),
  change_nameh: new rnoadm.gfx.Text('change name', '#fff', false),
  change_skin_color: new rnoadm.gfx.Text('change skin color', '#aaa', false),
  change_skin_colorh: new rnoadm.gfx.Text('change skin color', '#fff', false),
  change_shirt_color: new rnoadm.gfx.Text('change shirt color', '#aaa', false),
  change_shirt_colorh: new rnoadm.gfx.Text('change shirt color', '#fff', false),
  change_pants_color: new rnoadm.gfx.Text('change pants color', '#aaa', false),
  change_pants_colorh: new rnoadm.gfx.Text('change pants color', '#fff', false)
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

  return new rnoadm.hud.HUD('cc', function(w, h) {
    w /= rnoadm.gfx.TILE_SIZE * 2;
    h /= rnoadm.gfx.TILE_SIZE * 2;
    sprites.forEach(function(sprite) {
      sprite.paint(w - 3, h + 2);
    });
    name.paint(w, h + 3);
    rnoadm.hud.text_.change_name.paint(w, h + 3.5);
    change_gender.paint(w + 2, h - 2);
    rnoadm.hud.text_.change_skin_color.paint(w + 2, h - 1);
    rnoadm.hud.text_.change_shirt_color.paint(w + 2, h);
    rnoadm.hud.text_.change_pants_color.paint(w + 2, h + 1);
  });
};


rnoadm.net.addHandler('HUD', function(hud) {
  rnoadm.hud.activeHuds_.length = 0;
  if (hud['N'].length) {
    rnoadm.hud.activeHuds_.push(rnoadm.hud.huds_[hud['N']](hud['D'] || {}));
  }
  rnoadm.gfx.repaint();
});

rnoadm.gfx.paintHuds = function(w, h) {
  rnoadm.hud.activeHuds_.forEach(function(hud) {
    hud.paint_(w, h);
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
