goog.provide('rnoadm.net');

goog.require('goog.debug.Logger');
goog.require('goog.events');
goog.require('goog.net.WebSocket');
goog.require('goog.net.WebSocket.MessageEvent');


/**
 * Logger for rnoadm.net
 *
 * @type {goog.debug.Logger}
 * @private
 * @const
 */
rnoadm.net.logger_ = goog.debug.Logger.getLogger('rnoadm.net');


/**
 * The socket itself.
 *
 * @type {goog.net.WebSocket}
 * @private
 * @const
 */
rnoadm.net.socket_ = new goog.net.WebSocket();


/**
 *
 * @type {boolean}
 * @private
 */
rnoadm.net.connected_ = false;


/**
 * Sends a message to the server, or does nothing if the connection is down.
 *
 * @param {!Object} message the message to send.
 */
rnoadm.net.send = function(message) {
  if (rnoadm.net.connected_) {
    rnoadm.net.socket_.send(JSON.stringify(message));
  }
};


rnoadm.net.socket_.open('ws://' + location.host + '/ws');
goog.events.listen(rnoadm.net.socket_,
    goog.net.WebSocket.EventType.OPENED,
    rnoadm.net.onopen_);
goog.events.listen(rnoadm.net.socket_,
    goog.net.WebSocket.EventType.CLOSED,
    rnoadm.net.onclose_);
goog.events.listen(rnoadm.net.socket_,
    goog.net.WebSocket.EventType.MESSAGE,
    rnoadm.net.onmessage_);


/**
 * Called when the socket is opened.
 *
 * @private
 */
rnoadm.net.onopen_ = function() {
  rnoadm.net.connected_ = true;
};


/**
 * Called when the socket is closed.
 *
 * @private
 */
rnoadm.net.onclose_ = function() {
  rnoadm.net.connected_ = false;
};


/**
 * A 1-to-1 map of name => handler.
 *
 * @type {Object.<string, Function>}
 * @private
 * @const
 */
rnoadm.net.handlers_ = {};


/**
 * Adds a handler for a specific packet type.
 *
 * @param {string} name the name of the packet type.
 * @param {Function} fn the function to handle the packet type.
 */
rnoadm.net.addHandler = function(name, fn) {
  rnoadm.net.handlers_[name] = fn;
};


/**
 * Called when the socket recieves a message.
 *
 * @param {goog.net.WebSocket.MessageEvent} e an event containing the message.
 * @private
 */
rnoadm.net.onmessage_ = function(e) {
  var msg = JSON.parse(e.message), handler, name;
  for (name in msg) {
    if (handler = rnoadm.net.handlers_[name]) {
      handler(msg[name]);
    } else {
      rnoadm.net.logger_.info('Unhandled: ' + name);
    }
  }
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
