goog.provide('rnoadm.gfx.Text');

goog.require('rnoadm.gfx');



/**
 * @param {string} text
 * @param {string} color
 * @param {boolean} title
 * @param {boolean=} opt_left
 * @constructor
 * @struct
 */
rnoadm.gfx.Text = function(text, color, title, opt_left) {
  /**
   * @type {string}
   * @private
   * @const
   */
  this.text_ = text;

  /**
   * @type {string}
   * @private
   * @const
   */
  this.color_ = color;

  /**
   * @type {number}
   * @private
   * @const
   */
  this.height_ = title ? 1 : 0.5;

  /**
   * @type {string}
   * @private
   * @const
   */
  this.font_ = (this.height_ * rnoadm.gfx.TILE_SIZE) + 'px ' + (title ?
      '"Jolly Lodger"' : '"Open Sans Condensed"');

  /**
   * @type {boolean}
   * @private
   * @const
   */
  this.center_ = !opt_left;
};


rnoadm.gfx.Text.prototype.width = function() {
  rnoadm.gfx.ctx.font = this.font_;
  return rnoadm.gfx.ctx.measureText(this.text_).width / rnoadm.gfx.TILE_SIZE;
};


rnoadm.gfx.Text.prototype.height = function() {
  return this.height_;
};


rnoadm.gfx.Text.prototype.paint = function(x, y) {
  rnoadm.gfx.ctx.font = this.font_;
  var w = rnoadm.gfx.ctx.measureText(this.text_).width;
  x = Math.floor(x * rnoadm.gfx.TILE_SIZE - (this.center_ ? w / 2 : 0));
  y = Math.floor(y * rnoadm.gfx.TILE_SIZE);

  rnoadm.gfx.ctx.fillStyle = 'rgba(0,0,0,.2)';
  for (var x_ = -1; x_ <= 1; x_++) {
    for (var y_ = -1; y_ <= 2; y_++) {
      rnoadm.gfx.ctx.fillText(this.text_, x + x_, y + y_);
    }
  }

  rnoadm.gfx.ctx.fillStyle = this.color_;
  rnoadm.gfx.ctx.fillText(this.text_, x, y);
};


/**
 * @type {rnoadm.gfx.Text}
 * @const
 */
rnoadm.gfx.lostConnection = new rnoadm.gfx.Text('connection lost',
                                                '#aaa', true);

// vim: set tabstop=2 shiftwidth=2 expandtab:
