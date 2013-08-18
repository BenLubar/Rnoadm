goog.provide('rnoadm.main');

goog.require('goog.debug.FancyWindow');
goog.require('rnoadm.login');
goog.require('rnoadm.net');

new goog.debug.FancyWindow().init();


/**
 * The hash of the compiled client code, as sent by the server.
 *
 * @type {string}
 * @private
 */
rnoadm.main.clientHash_;


rnoadm.net.addHandler('ClientHash', function(hash) {
  if (goog.isDef(rnoadm.main.clientHash_)) {
    if (rnoadm.main.clientHash_ != hash) {
      location.reload(true);
    }
  } else {
    rnoadm.main.clientHash_ = hash;
  }
});


// vim: set tabstop=2 shiftwidth=2 expandtab:
