goog.provide('rnoadm.gfx.NetworkSprite');
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
  var scale = Math.floor(opt_scale) || 1;

  /**
   * @type {number}
   * @const
   * @private
   */
  this.scale_ = scale;

  /** @type {Image} */
  var img;

  function init() {
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

    var noRecolor = color == 'no';
    if (noRecolor) {
      color = '#000';
    }

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
      if (noRecolor) {
        return x;
      }
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
    for (var sy = 0; sy < canvas.height; sy++) {
      for (var sx = 0; sx < canvas.width; sx++) {
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
  }

  image = image + '.png';
  if (img = rnoadm.gfx.Sprite.cache_[image]) {
    init();
  } else {
    img = new Image();
    img.onload = function() {
      rnoadm.gfx.Sprite.cache_[image] = img;
      init();
    };
    img.src = image;
  }
};


/**
 * @type {Object.<string, Image>}
 * @private
 * @const
 */
rnoadm.gfx.Sprite.cache_ = {};


/**
 * @param {number} x The x coordinate to draw on the canvas, in tiles.
 * @param {number} y The y coordinate to draw on the canvas, in tiles.
 */
rnoadm.gfx.Sprite.prototype.paint = function(x, y) {
  /** @type {number} */
  var sx = this.xOffset_;
  /** @type {number} */
  var sy = this.yOffset_;

  var sprite = this;

  function draw() {
    rnoadm.gfx.ctx.drawImage(sprite.canvas_,
        Math.floor(sx * sprite.width_ * sprite.scale_),
        Math.floor(sy * sprite.height_ * sprite.scale_),
        Math.floor(sprite.width_ * sprite.scale_),
        Math.floor(sprite.height_ * sprite.scale_),
        Math.floor(x * rnoadm.gfx.TILE_SIZE -
        (sprite.width_ * sprite.scale_ - rnoadm.gfx.TILE_SIZE) / 2),
        Math.floor(y * rnoadm.gfx.TILE_SIZE - sprite.height_ * sprite.scale_),
        Math.floor(sprite.width_ * sprite.scale_),
        Math.floor(sprite.height_ * sprite.scale_));
  }

  if (rnoadm.gfx.Sprite.floorPass) {
    if (this.animation_ == '_fl') {
      y += 0.5;
      draw();
    }
    return;
  }

  switch (this.animation_) {
    case '':
      break;
    case 'ccr': // character creation rotation
      sx += [0, 6, 3, 9][Math.floor(Date.now() / 1500) % 4];
      rnoadm.gfx.repaint(1500 - Date.now() % 1500);
      break;
    case 'wa': // walk (alternating)
      sx += [0, 1, 0, 2][Math.floor(Date.now() / 150) % 4];
      rnoadm.gfx.repaint(150 - Date.now() % 150);
      break;
    case 'l2': // loop (three frames)
      sx += Math.floor(Date.now() / 150) % 2;
      rnoadm.gfx.repaint(150 - Date.now() % 150);
      break;
    case 'l3': // loop (three frames)
      sx += Math.floor(Date.now() / 150) % 3;
      rnoadm.gfx.repaint(150 - Date.now() % 150);
      break;
    case '_ac': // special (admin crown)
      if (sy < 2)
        break;
      var time = Date.now() / 10000;
      switch (sy) {
        case 2:
        case 3:
        case 4:
        case 5:
          y += Math.sin(time * 5 + Math.cos(time * 3) * 7 + sy) / 8;
          break;
      }
      rnoadm.gfx.repaint(100);
      break;
    case 'wa_ac': // walk (alternating) special (admin crown)
      sx += [0, 1, 0, 2][Math.floor(Date.now() / 150) % 4];
      rnoadm.gfx.repaint(50);
      if (sy < 2)
        break;
      var time = Date.now() / 10000;
      switch (sy) {
        case 2:
        case 3:
        case 4:
        case 5:
          y += Math.sin(time * 5 + Math.cos(time * 3) * 7 + sy) / 8;
          break;
      }
      break;
    case '_fl': // special (floor)
      return;
  }

  draw();
};


/**
 * @type {boolean}
 */
rnoadm.gfx.Sprite.floorPass = false;


/**
 * @return {boolean}
 */
rnoadm.gfx.Sprite.prototype.isFloor = function() {
  return this.animation_ == '_fl';
};


/**
 * @param {rnoadm.gfx.NetworkSprite} sprite
 * @return {rnoadm.gfx.Sprite}
 */
rnoadm.gfx.Sprite.fromNetwork = function(sprite) {
  return new rnoadm.gfx.Sprite(sprite['S'], sprite['C'],
      sprite['E']['a'] || '', sprite['E']['x'] || 0, sprite['E']['y'] || 0,
      sprite['E']['w'] || rnoadm.gfx.TILE_SIZE,
      sprite['E']['h'] || rnoadm.gfx.TILE_SIZE,
      sprite['E']['s'] || 1);
};


/**
 * @typedef {{S:string, C:string, E:{a:string, w:number, h:number,
 *                                   x:number, y:number, s:number}}}
 */
rnoadm.gfx.NetworkSprite;

// vim: set tabstop=2 shiftwidth=2 expandtab:
