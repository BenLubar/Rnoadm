goog.provide('rnoadm.state.NetworkObject');
goog.provide('rnoadm.state.Object');
goog.provide('rnoadm.state.lerp');

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

  (net['S'] || []).forEach(function(sprite) {
    sprites.push(rnoadm.gfx.Sprite.fromNetwork(sprite));
  });

  /** @type {Array.<string>} */
  this.actions = net['A'] || [];
};


/**
 * @return {rnoadm.state.Object}
 */
rnoadm.state.Object.prototype.clone = function() {
  var copy = new rnoadm.state.Object(this.x, this.y, this.id, {
    'N': this.name,
    'A': this.actions,
    'S': null
  });
  copy.sprites = this.sprites;
  return copy;
};


/**
 * @param {!rnoadm.state.NetworkObject} net The object recieved by rnoadm.net.
 */
rnoadm.state.Object.prototype.update = function(net) {
  /** @type {string} */
  this.name = net['N'];

  /** @type {Array.<string>} */
  this.actions = net['A'] || [];

  /** @type {Array.<rnoadm.gfx.Sprite>} */
  var sprites = this.sprites;

  sprites.length = 0;

  (net['S'] || []).forEach(function(sprite) {
    sprites.push(rnoadm.gfx.Sprite.fromNetwork(sprite));
  });
};


/**
 * @param {number} x The current X coordinate of the object.
 * @param {number} y The current Y coordinate of the object.
 */
rnoadm.state.Object.prototype.move = function(x, y) {
  this.prevX = rnoadm.state.lerp(this.x, this.prevX, this.lastChange);
  this.prevY = rnoadm.state.lerp(this.y, this.prevY, this.lastChange);
  this.x = x;
  this.y = y;
  this.lastChange = Date.now();
};


/**
 * @param {number} xOffset
 * @param {number} yOffset
 */
rnoadm.state.Object.prototype.paint = function(xOffset, yOffset) {
  /** @type {number} */
  var x = rnoadm.state.lerp(this.x, this.prevX, this.lastChange);
  /** @type {number} */
  var y = rnoadm.state.lerp(this.y, this.prevY, this.lastChange);

  this.sprites.forEach(function(sprite) {
    sprite.paint(x + xOffset, y + yOffset);
  });
};


/** @typedef {{N:string, S:Array.<rnoadm.gfx.NetworkSprite>,
 *             A:Array.<string>}} */
rnoadm.state.NetworkObject;


/**
 * @param {number} current
 * @param {number} prev
 * @param {number} t
 * @return {number}
 */
rnoadm.state.lerp = function(current, prev, t) {
  var time = Date.now() - t;
  if (time < 400) {
    rnoadm.gfx.repaint();
    return (time * current + (400 - time) * prev) / 400;
  }
  return current;
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
