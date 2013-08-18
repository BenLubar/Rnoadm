goog.provide('rnoadm.gfx.Sprite');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');



/**
 * @constructor
 * @struct
 * @param {string} image The name of the sprite sheet to use.
 * @param {string} color A CSS color to tint the image.
 * @param {string} animation The animation function to use.
 * @param {number} xOffset The number of sprites to skip from the left.
 * @param {number} yOffset The number of sprites to skip from the top.
 * @param {number} width The width of a sprite in pixels.
 * @param {number} height The height of a sprite in pixels.
 * @param {number=} opt_scale An integer multiplier for the sprite's size.
 */
rnoadm.gfx.Sprite = function(image, color, animation, xOffset, yOffset,
                             width, height, opt_scale) {
  /**
   * @type {HTMLCanvasElement}
   * @const
   */
  var canvas = goog.asserts.assertInstanceof(document.createElement('canvas'),
                                             HTMLCanvasElement);

  /**
   * @type {HTMLCanvasElement}
   * @const
   * @private
   */
  this.canvas_ = canvas;

  /**
   * @type {string}
   * @const
   * @private
   */
  this.animation_ = animation;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.xOffset_ = xOffset;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.yOffset_ = yOffset;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.width_ = width;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.height_ = height;

  /**
   * @type {number}
   * @const
   */
  var scale = Math.floor(scale) || 1;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.scale_ = scale;

  /**
   * @type {Image}
   * @const
   */
  var img = new Image();
  img.onload = function() {
    if (!img.width || !img.height) {
      return;
    }
    canvas.width = img.width;
    canvas.height = img.height;
    /**
     * @type {CanvasRenderingContext2D}
     * @const
     */
    var ctx = goog.asserts.assertInstanceof(canvas.getContext('2d'),
                                            CanvasRenderingContext2D);

    // Determine the actual color from the string color.
    ctx.fillStyle = color;
    ctx.fillRect(0, 0, 1, 1);
    /** @type {ImageData} */
    var pix = ctx.getImageData(0, 0, 1, 1);
    /** @type {number} */
    var r = pix.data[0];
    /** @type {number} */
    var g = pix.data[1];
    /** @type {number} */
    var b = pix.data[2];
    /** @type {number} */
    var a = pix.data[3];
    ctx.clearRect(0, 0, 1, 1);

    // Draw the base image.
    ctx.drawImage(img, 0, 0);
    pix = ctx.getImageData(0, 0, img.width, img.height);

    // Scale the canvas. This also clears it.
    canvas.width *= scale;
    canvas.height *= scale;

    /**
     * @param {number} x base color
     * @param {number} y accent color
     * @return {number}
     */
    function fade(x, y) {
      if (x >= 128) {
        return 255 - (255 - x) * (255 - y) / 127;
      }
      return x * y / 127;
    }

    // Here comes the fun part!
    /** @type {ImageData} */
    var scaled = ctx.getImageData(0, 0, canvas.width, canvas.height);
    var rowIndex = 0;
    var baseIndex = 0;
    var scaledIndex = 0;
    for (var sy = 0; sy < canvas.width; sy++) {
      for (var sx = 0; sx < canvas.height; sx++) {
        scaled.data[scaledIndex + 0] = fade(pix.data[baseIndex + 0], r);
        scaled.data[scaledIndex + 1] = fade(pix.data[baseIndex + 1], g);
        scaled.data[scaledIndex + 2] = fade(pix.data[baseIndex + 2], b);
        scaled.data[scaledIndex + 3] = pix.data[baseIndex + 3] * a / 255;
        if (sx % scale == scale - 1) {
          baseIndex += 4;
        }
        scaledIndex += 4;
      }
      if (sy % scale == scale - 1) {
        rowIndex = baseIndex;
      } else {
        baseIndex = rowIndex;
      }
    }
    ctx.putImageData(scaled, 0, 0);

    rnoadm.gfx.repaint();
  };
  img.src = image + '.png?' + rnoadm.gfx.GRAPHICS_REVISION;
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
