goog.provide('rnoadm.state.NetworkObject');
goog.provide('rnoadm.state.Object');

goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.NetworkSprite');
goog.require('rnoadm.gfx.Sprite');



/**
 * @param {number} x The current X coordinate of the object.
 * @param {number} y The current Y coordinate of the object.
 * @param {string} id The network ID of the object.
 * @param {!rnoadm.state.NetworkObject} net The object recieved by rnoadm.net.
 * @constructor
 * @struct
 */
rnoadm.state.Object = function(x, y, id, net) {
  /**
   * @type {string}
   * @const
   */
  this.id = id;

  /** @type {number} */
  this.x = x;
  /** @type {number} */
  this.prevX = x;

  /** @type {number} */
  this.y = y;
  /** @type {number} */
  this.prevY = y;

  /** @type {number} */
  this.lastChange = 0;

  /** @type {string} */
  this.name = net['N'];

  /** @type {Array.<rnoadm.gfx.Sprite>} */
  var sprites = [];

  /** @type {Array.<rnoadm.gfx.Sprite>} */
  this.sprites = sprites;

  net['S'].forEach(function(sprite) {
    sprites.push(rnoadm.gfx.Sprite.fromNetwork(sprite));
  });
};


/**
 * @param {!rnoadm.state.NetworkObject} net The object recieved by rnoadm.net.
 */
rnoadm.state.Object.prototype.update = function(net) {
  this.name = net['N'];

  /** @type {Array.<rnoadm.gfx.Sprite>} */
  var sprites = this.sprites;

  sprites.length = 0;

  net['S'].forEach(function(sprite) {
    sprites.push(new rnoadm.gfx.Sprite(sprite['S'], sprite['C'],
        sprite['E']['a'], sprite['E']['x'], sprite['E']['y'],
        sprite['E']['w'] || rnoadm.gfx.TILE_SIZE,
        sprite['E']['h'] || rnoadm.gfx.TILE_SIZE,
        sprite['E']['s'] || 1));
  });
};


/**
 * @param {number} x The current X coordinate of the object.
 * @param {number} y The current Y coordinate of the object.
 */
rnoadm.state.Object.prototype.move = function(x, y) {
  this.lastChange = Date.now();
  this.prevX = this.x;
  this.prevY = this.y;
  this.x = x;
  this.y = y;
};


/**
 * @param {number} xOffset
 * @param {number} yOffset
 */
rnoadm.state.Object.prototype.paint = function(xOffset, yOffset) {
  /** @type {number} */
  var x = this.x;
  /** @type {number} */
  var y = this.y;

  /** @type {number} */
  var time = Date.now() - this.lastChange;
  if (time < 400) {
    x = (time * x + (400 - time) * this.prevX) / 400;
    y = (time * y + (400 - time) * this.prevY) / 400;
    rnoadm.gfx.repaint();
  }

  this.sprites.forEach(function(sprite) {
    sprite.paint(x + xOffset, y + yOffset);
  });
};


/** @typedef {{N:string, S:Array.<rnoadm.gfx.NetworkSprite>}} */
rnoadm.state.NetworkObject;

// vim: set tabstop=2 shiftwidth=2 expandtab:
