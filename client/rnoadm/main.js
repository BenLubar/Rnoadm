goog.provide('rnoadm.main');

goog.require('goog.debug.FancyWindow');
goog.require('goog.debug.Logger');
goog.require('rnoadm.hud');
goog.require('rnoadm.login');
goog.require('rnoadm.net');
goog.require('rnoadm.state');


if (goog.DEBUG) {
  new goog.debug.FancyWindow();
}


/**
 * Logger for rnoadm.main
 *
 * @type {goog.debug.Logger}
 * @private
 * @const
 */
rnoadm.main.logger_ = goog.debug.Logger.getLogger('rnoadm.main');


/**
 * The hash of the compiled client code, as sent by the server.
 *
 * @type {string|undefined}
 * @private
 */
rnoadm.main.clientHash_;


rnoadm.net.addHandler('ClientHash', function(hash) {
  rnoadm.main.logger_.info('Client hash: ' + hash);
  if (goog.isDef(rnoadm.main.clientHash_)) {
    if (rnoadm.main.clientHash_ != hash) {
      location.reload(true);
    }
  } else {
    rnoadm.main.clientHash_ = hash;
  }
});


// vim: set tabstop=2 shiftwidth=2 expandtab:
