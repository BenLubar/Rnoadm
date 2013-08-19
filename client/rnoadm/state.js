goog.provide('rnoadm.state');

goog.require('goog.asserts');
goog.require('goog.debug.Logger');
goog.require('rnoadm.gfx');
goog.require('rnoadm.gfx.NetworkSprite');
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
  if (time < 500) {
    x = (time * x + (500 - time) * this.prevX) / 500;
    y = (time * y + (500 - time) * this.prevY) / 500;
    rnoadm.gfx.repaint();
  }

  this.sprites.forEach(function(sprite) {
    sprite.paint(x + xOffset, y + yOffset);
  });
};


/** @typedef {{N:string, S:Array.<rnoadm.gfx.NetworkSprite>}} */
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
      if (id in rnoadm.state.objects_[index]) {
        rnoadm.state.objects_[index][id].update(object);
      } else {
        rnoadm.state.objects_[index][id] = new rnoadm.state.Object(x, y, id,
                                                                   object);
      }
      return;
    }

    goog.asserts.assert(goog.isDef(fromX) && goog.isDef(fromY),
                        'update: blank update');
    var fromIndex = fromX | fromY << 8;
    goog.asserts.assert(id in rnoadm.state.objects_[fromIndex],
                        'update: move of nonexistent object');
    rnoadm.state.objects_[index][id] = rnoadm.state.objects_[fromIndex][id];
    delete rnoadm.state.objects_[fromIndex][id];
    rnoadm.state.objects_[index][id].move(x, y);
  });
  rnoadm.gfx.repaint();
});


/**
 * @type {number}
 * @private
 */
rnoadm.state.playerLastMove_ = 0;


/**
 * @type {number}
 * @private
 */
rnoadm.state.playerPrevX_ = 127;


/**
 * @type {number}
 * @private
 */
rnoadm.state.playerX_ = 127;


rnoadm.net.addHandler('PlayerX', function(x) {
  rnoadm.state.playerPrevX_ = rnoadm.state.playerX_;
  rnoadm.state.playerX_ = x;
  rnoadm.state.playerLastMove_ = Date.now();
  rnoadm.gfx.repaint();
});


/**
 * @type {number}
 * @private
 */
rnoadm.state.playerPrevY_ = 127;


/**
 * @type {number}
 * @private
 */
rnoadm.state.playerY_ = 127;


rnoadm.net.addHandler('PlayerY', function(y) {
  rnoadm.state.playerPrevY_ = rnoadm.state.playerY_;
  rnoadm.state.playerY_ = y;
  rnoadm.state.playerLastMove_ = Date.now();
  rnoadm.gfx.repaint();
});


rnoadm.net.onConnect.push(function() {
  for (var i = 0; i < rnoadm.state.objects_.length; i++) {
    delete rnoadm.state.objects_[i];
  }
});


/**
 * @type {rnoadm.gfx.Sprite}
 * @private
 * @const
 */
rnoadm.state.grass_ = new rnoadm.gfx.Sprite('grass', '#268f1e', '',
                                            0, 0, 512, 512);


rnoadm.gfx.paintObjects = function(w, h) {
  /** @type {number} */
  var playerX = rnoadm.state.playerX_;
  /** @type {number} */
  var playerY = rnoadm.state.playerY_;

  /** @type {number} */
  var time = Date.now() - rnoadm.state.playerLastMove_;
  if (time < 500) {
    playerX = (time * playerX + (500 - time) * rnoadm.state.playerPrevX_) / 500;
    playerY = (time * playerY + (500 - time) * rnoadm.state.playerPrevY_) / 500;
    rnoadm.gfx.repaint();
  }

  var xOffset = w / 2 / rnoadm.gfx.TILE_SIZE - playerX;
  var yOffset = h / 2 / rnoadm.gfx.TILE_SIZE - playerY;
  for (var x = 0; x < 256; x += 512 / rnoadm.gfx.TILE_SIZE) {
    for (var y = 0; y < 256; y += 512 / rnoadm.gfx.TILE_SIZE) {
      rnoadm.state.grass_.paint(xOffset + x + 512 / rnoadm.gfx.TILE_SIZE / 2,
                                yOffset + y + 512 / rnoadm.gfx.TILE_SIZE / 2);
    }
  }
  for (var x = Math.max(0, Math.floor(xOffset - w / 2));
       x < Math.min(256, Math.floor(xOffset + w / 2));
       x++) {
    for (var y = Math.max(0, Math.floor(yOffset - h / 2));
         y < Math.min(256, Math.floor(yOffset + h / 2));
         y++) {
      var objects = rnoadm.state.objects_[x | y << 8];
      if (objects) {
        for (var id in objects) {
          objects[id].paint(xOffset, yOffset);
        }
      }
    }
  }
};

rnoadm.gfx.clickObject = function(x, y, w, h) {
  var xOffset = w / 2 / rnoadm.gfx.TILE_SIZE - rnoadm.state.playerX_;
  var yOffset = h / 2 / rnoadm.gfx.TILE_SIZE - rnoadm.state.playerY_;
  x = Math.floor(x / rnoadm.gfx.TILE_SIZE - xOffset);
  y = Math.ceil(y / rnoadm.gfx.TILE_SIZE - yOffset);
  if (x < 0 || x >= 256 || y < 0 || y >= 256)
    return false;
  rnoadm.net.send({'Walk': {'X': x, 'Y': y}});
  return true;
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
