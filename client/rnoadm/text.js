goog.provide('rnoadm.gfx.Text');

goog.require('rnoadm.gfx');


/**
 * @constructor
 * @struct
 */
rnoadm.gfx.Text = function(text, color, title) {
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
   * @type {string}
   * @private
   * @const
   */
  this.font_ = title ? '32px "Jolly Lodger"' : '16px "Open Sans Condensed"';
};


rnoadm.gfx.Text.prototype.paint = function(x, y) {
  rnoadm.gfx.ctx.font = this.font_;
  var w = rnoadm.gfx.ctx.measureText(this.text_).width;
  x = Math.floor(x * rnoadm.gfx.TILE_SIZE - w / 2);
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
