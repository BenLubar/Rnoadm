goog.provide('rnoadm.gfx');

goog.require('goog.asserts');
goog.require('rnoadm.net');


/** @define {number} The dimension, in pixels, of a "tile". */
rnoadm.gfx.TILE_SIZE = 32;


/** @define {number} The number of times a sprite sheet has changed. */
rnoadm.gfx.GRAPHICS_REVISION = 13;


/**
 * @type {HTMLCanvasElement}
 * @const
 */
rnoadm.gfx.canvas = goog.asserts.assertInstanceof(
    document.createElement('canvas'),
    HTMLCanvasElement);


/**
 * @type {CanvasRenderingContext2D}
 * @const
 */
rnoadm.gfx.ctx = goog.asserts.assertInstanceof(
    rnoadm.gfx.canvas.getContext('2d'),
    CanvasRenderingContext2D);


/**
 * @type {number}
 * @private
 */
rnoadm.gfx.width_ = window.innerWidth;


/**
 * @type {number}
 * @private
 */
rnoadm.gfx.height_ = window.innerHeight;


/**
 * @type {number}
 */
rnoadm.gfx.canvas.width = rnoadm.gfx.width_;


/**
 * @type {number}
 */
rnoadm.gfx.canvas.height = rnoadm.gfx.height_;

window.addEventListener('resize', function() {
  rnoadm.gfx.canvas.width = rnoadm.gfx.width_ = window.innerWidth;
  rnoadm.gfx.canvas.height = rnoadm.gfx.height_ = window.innerHeight;
  rnoadm.gfx.repaint();
}, false);


/**
 * @type {number}
 * @private
 */
rnoadm.gfx.nextRepaint_ = Infinity;


/**
 * @type {number}
 * @private
 */
rnoadm.gfx.repaintTimeout_ = 0;


/**
 * Schedules a repaint. Does not actually perform the repaint.
 *
 * @param {number=} opt_delay milliseconds to wait before repainting.
 */
rnoadm.gfx.repaint = function(opt_delay) {
  var delay = !isNaN(opt_delay) && opt_delay > 0 ? opt_delay : 0;
  var time = delay + Date.now();
  if (rnoadm.gfx.nextRepaint_ > time) {
    if (rnoadm.gfx.repaintTimeout_) {
      clearTimeout(rnoadm.gfx.repaintTimeout_);
    }
    rnoadm.gfx.repaintTimeout_ = setTimeout(function() {
      window.requestAnimationFrame(rnoadm.gfx.paint_);
      rnoadm.gfx.repaintTimeout_ = 0;
      rnoadm.gfx.nextRepaint_ = Infinity;
    }, delay);
    rnoadm.gfx.nextRepaint_ = time;
  }
};


/**
 * @type {boolean}
 * @private
 */
rnoadm.gfx.connected_ = false;

rnoadm.net.onConnect.push(function() {
  rnoadm.gfx.connected_ = true;
  rnoadm.gfx.repaint();
});

rnoadm.net.onDisconnect.push(function() {
  rnoadm.gfx.connected_ = false;
  rnoadm.gfx.repaint();
});


/**
 * @private
 */
rnoadm.gfx.paint_ = function() {
  rnoadm.gfx.ctx.clearRect(0, 0, rnoadm.gfx.width_, rnoadm.gfx.height_);
  if (rnoadm.gfx.connected_) {
    rnoadm.gfx.paintObjects(rnoadm.gfx.width_, rnoadm.gfx.height_);
    rnoadm.gfx.paintHuds(rnoadm.gfx.width_, rnoadm.gfx.height_);
  } else {
    rnoadm.gfx.lostConnection.paint(
        rnoadm.gfx.width_ / 2 / rnoadm.gfx.TILE_SIZE,
        rnoadm.gfx.height_ / 2 / rnoadm.gfx.TILE_SIZE);
  }
};


/** @type {function(number, number)} */
rnoadm.gfx.paintObjects;


/** @type {function(number, number)} */
rnoadm.gfx.paintHuds;


/** @type {function(number, number, number, number):boolean} */
rnoadm.gfx.mouseMovedHud;


/** @type {function(number, number, number, number):boolean} */
rnoadm.gfx.clickHud;


/** @type {function(number, number, number, number):boolean} */
rnoadm.gfx.clickObject;


/** @type {function(number, number, number, number):boolean} */
rnoadm.gfx.rightClickHud;


rnoadm.gfx.canvas.onmousemove = function(e) {
  if (rnoadm.gfx.mouseMovedHud(e.offsetX, e.offsetY,
                               rnoadm.gfx.width_,
                               rnoadm.gfx.height_))
    return;
};


rnoadm.gfx.canvas.onmouseout = function() {
  if (rnoadm.gfx.mouseMovedHud(-Infinity, -Infinity,
                               rnoadm.gfx.width_,
                               rnoadm.gfx.height_))
    return;
};


/**
 * @type {boolean}
 * @private
 */
rnoadm.gfx.focused_ = true;


window.addEventListener('blur', function() {
  rnoadm.gfx.focused_ = false;
}, false);

window.addEventListener('focus', function() {
  setTimeout(function() {
    rnoadm.gfx.focused_ = true;
  }, 100);
}, false);


rnoadm.gfx.canvas.onclick = function(e) {
  if (!rnoadm.gfx.focused_) {
    rnoadm.gfx.focused_ = true;
    return;
  }
  if (rnoadm.gfx.clickHud(e.offsetX, e.offsetY,
                          rnoadm.gfx.width_,
                          rnoadm.gfx.height_))
    return;
  if (rnoadm.gfx.clickObject(e.offsetX, e.offsetY,
                             rnoadm.gfx.width_,
                             rnoadm.gfx.height_))
    return;
};


rnoadm.gfx.canvas.oncontextmenu = function(e) {
  e.preventDefault();
  if (rnoadm.gfx.rightClickHud(e.offsetX, e.offsetY,
                               rnoadm.gfx.width_,
                               rnoadm.gfx.height_))
    return;
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
