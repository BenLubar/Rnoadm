goog.provide('rnoadm.state.text');

goog.require('rnoadm.gfx.Text');
goog.require('rnoadm.net');
goog.require('rnoadm.state.lerp');



/**
 * @constructor
 * @struct
 */
rnoadm.state.text.Text = function(x, y, text, color) {
  this.x_ = x;
  this.y_ = y;
  this.start_ = goog.now();
  this.text_ = new rnoadm.gfx.Text(text, color, false);
};


rnoadm.state.text.Text.prototype.paint = function(xOffset, yOffset) {
  var y = rnoadm.state.lerp(this.y_ - 0.8, this.y_ - 0.2, this.start_) + yOffset;
  var opacity = Math.min(rnoadm.state.lerp(0, 2, this.start_), 1);

  rnoadm.gfx.ctx.globalAlpha = opacity;
  this.text_.paint(this.x_ + 0.5 + xOffset, y);
  rnoadm.gfx.ctx.globalAlpha = 1;

  return opacity > 0;
};


/**
 * @type {Array.<rnoadm.state.text.Text>}
 * @private
 * @const
 */
rnoadm.state.text.text_ = [];


rnoadm.net.addHandler('Ftxt', function(text) {
  rnoadm.state.text.text_.push(new rnoadm.state.text.Text(
      text['X'], text['Y'], text['T'], text['C']));
  rnoadm.gfx.repaint();
});


rnoadm.state.text.paint = function(xOffset, yOffset) {
  if (rnoadm.state.text.text_.length)
    rnoadm.gfx.repaint(20);
  rnoadm.state.text.text_.filter(function(text) {
    return text.paint(xOffset, yOffset);
  });
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
