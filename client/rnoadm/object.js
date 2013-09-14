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

  /** @type {(Array.<number>|undefined)} */
  this.health = net['H'];
};


/**
 * @return {rnoadm.state.Object}
 */
rnoadm.state.Object.prototype.clone = function() {
  var copy = new rnoadm.state.Object(this.x, this.y, this.id, {
    'N': this.name,
    'A': this.actions,
    'S': null,
    'H': this.health
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

  /** @type {(Array.<number>|undefined)} */
  this.health = net['H'];
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
  this.lastChange = goog.now();
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

  var maxWidth = 0, maxHeight = 0;
  this.sprites.forEach(function(sprite) {
    maxWidth = Math.max(maxWidth, sprite.width());
    maxHeight = Math.max(maxHeight, sprite.height());
    sprite.paint(x + xOffset, y + yOffset);
  });

  if (this.health && this.sprites.length) {
    var healthX = (x + xOffset) * rnoadm.gfx.TILE_SIZE;
    var healthY = (y + yOffset) * rnoadm.gfx.TILE_SIZE - maxHeight - 5;
    var h0 = this.health[0], h1 = this.health[1];
    rnoadm.gfx.ctx.fillStyle = '#000';
    rnoadm.gfx.ctx.fillRect(healthX, healthY, maxWidth, 3);
    rnoadm.gfx.ctx.fillStyle = 'rgb(' +
        Math.floor(Math.min(Math.max(510 - h0 / h1 * 510, 0), 255)) + ',' +
        Math.floor(Math.min(Math.max(h0 / h1 * 510, 0), 255)) + ',0)';
    rnoadm.gfx.ctx.fillRect(healthX, healthY, maxWidth * h0 / h1, 3);
  }
};


/** @typedef {{N:string, S:Array.<rnoadm.gfx.NetworkSprite>,
 *             A:Array.<string>, H:(Array.<number>|undefined)}} */
rnoadm.state.NetworkObject;


/**
 * @param {number} current
 * @param {number} prev
 * @param {number} t
 * @return {number}
 */
rnoadm.state.lerp = function(current, prev, t) {
  var time = goog.now() - t;
  if (time < 400) {
    rnoadm.gfx.repaint();
    return (time * current + (400 - time) * prev) / 400;
  }
  return current;
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
