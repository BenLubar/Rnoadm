goog.provide('rnoadm.state');

goog.require('goog.asserts');
goog.require('goog.debug.Logger');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.Sprite');
goog.require('rnoadm.net');


/**
 * Logger for rnoadm.state
 *
 * @type {goog.debug.Logger}
 * @private
 * @const
 */
rnoadm.state.logger_ = goog.debug.Logger.getLogger('rnoadm.state');


/**
 * @type {Array.<Object.<string, rnoadm.state.Object>>}
 * @private
 * @const
 */
rnoadm.state.objects_ = Array(256 * 256);



/**
 * @param {number} x The current X coordinate of the object.
 * @param {number} y The current Y coordinate of the object.
 * @param {string} id The network ID of the object.
 * @param {!rnoadm.state.NetworkObject} net The object recieved by rnoadm.net.
 * @constructor
 * @struct
 */
rnoadm.state.Object = function(x, y, id, net) {
  /** @type {number} */
  this.x = x;
  /** @type {number} */
  this.prevX = x;

  /** @type {number} */
  this.y = y;
  /** @type {number} */
  this.prevY = y;

  /** @type {number} */
  this.lastChange = Date.now();

  /** @type {string} */
  this.name = net['N'];

  /** @type {Array.<rnoadm.gfx.Sprite>} */
  this.sprites = [];

  net['S'].forEach(function(sprite) {
    this.sprites.push(new rnoadm.gfx.Sprite(sprite['S'], sprite['C'],
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
rnoadm.state.Object.prototype.update = function(x, y) {
  this.prevX = this.x;
  this.prevY = this.y;
  this.x = x;
  this.y = y;
};


/**
 * @typedef {{S:string, C:string, E:{a:string, w:number, h:number,
 *                                   x:number, y:number, s:number}}} */
rnoadm.state.NetworkSprite;


/** @typedef {{N:string, S:Array.<rnoadm.state.NetworkSprite>}} */
rnoadm.state.NetworkObject;

rnoadm.net.addHandler('Update', function(updates) {
  updates.forEach(function(update) {
    /**
     * @type {string}
     * @const
     */
    var id = update['I'];
    /**
     * @type {number}
     * @const
     */
    var x = update['X'];
    /**
     * @type {number}
     * @const
     */
    var y = update['Y'];
    /**
     * @type {number|undefined}
     * @const
     */
    var fromX = update['Fx'];
    /**
     * @type {number|undefined}
     * @const
     */
    var fromY = update['Fy'];
    /**
     * @type {boolean|undefined}
     * @const
     */
    var removed = update['R'];
    /**
     * @type {(!rnoadm.state.NetworkObject|undefined)}
     * @const
     */
    var object = update['O'];

    goog.asserts.assert(x >= 0 && x < 256, 'update: x out of range');
    goog.asserts.assert(y >= 0 && y < 256, 'update: y out of range');

    /**
     * @type {number}
     * @const
     */
    var index = x | y << 8;

    if (!rnoadm.state.objects_[index]) {
      rnoadm.state.objects_[index] = {};
    }

    if (removed) {
      goog.asserts.assert(!goog.isDef(fromX) && !goog.isDef(fromY) &&
                          !goog.isDef(object), 'update: contradictory update');
      goog.asserts.assert(id in rnoadm.state.objects_[index],
                          'update: removal of nonexistent object');
      delete rnoadm.state.objects_[index][id];
      rnoadm.gfx.repaint();
      return;
    }

    if (goog.isDef(object)) {
      goog.asserts.assert(!goog.isDef(fromX) && !goog.isDef(fromY),
                          'update: contradictory update');
      rnoadm.state.objects_[index][id] = new rnoadm.state.Object(x, y, id,
                                                                 object);
      return;
    }

    goog.asserts.assert(goog.isDef(fromX) && goog.isDef(fromY),
                        'update: blank update');
    var fromIndex = fromX | fromY << 8;
    goog.asserts.assert(id in rnoadm.state.objects_[fromIndex],
                        'update: move of nonexistent object');
    rnoadm.state.objects_[index][id] = rnoadm.state.objects_[fromIndex][id];
    delete rnoadm.state.objects_[fromIndex][id];
    rnoadm.state.objects_[index][id].update(x, y);
  });
});

// vim: set tabstop=2 shiftwidth=2 expandtab:
