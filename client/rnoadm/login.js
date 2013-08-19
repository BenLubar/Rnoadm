goog.provide('rnoadm.login');

goog.require('goog.asserts');
goog.require('rnoadm.gfx');
goog.require('rnoadm.net');


/**
 * @type {string}
 */
rnoadm.login.username = window.sessionStorage['rnoadm_username'] || '';


/**
 * @type {string}
 */
rnoadm.login.password = window.sessionStorage['rnoadm_password'] || '';


rnoadm.net.onConnect.push(function() {
  if (rnoadm.login.username.length && rnoadm.login.password.length > 2) {
    rnoadm.login.form_.submit();
  }
});


/**
 * @type {HTMLFormElement}
 * @private
 * @const
 */
rnoadm.login.form_ = goog.asserts.assertInstanceof(
    document.querySelector('form'),
    HTMLFormElement);


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.username_ = goog.asserts.assertInstanceof(
    rnoadm.login.form_.querySelector('#username'),
    HTMLInputElement);


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.password_ = goog.asserts.assertInstanceof(
    rnoadm.login.form_.querySelector('#password'),
    HTMLInputElement);


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.dummy_ = goog.asserts.assertInstanceof(
    rnoadm.login.form_.querySelector('#password2'),
    HTMLInputElement);


/**
 * @param {...string} var_args
 * @private
 */
rnoadm.login.admin_ = function(var_args) {
  rnoadm.net.send({'Admin': [].slice.call(arguments)});
};

goog.exportSymbol('admin', rnoadm.login.admin_);


/**
 * @private
 */
rnoadm.login.onlogin_ = function() {
  var parent = rnoadm.login.form_.parentNode;
  parent.removeChild(rnoadm.login.form_);
  parent.style.overflow = 'hidden';
  parent.style.fontSize = '0';
  parent.appendChild(rnoadm.gfx.canvas);
};


/**
 */
rnoadm.login.password_.onchange = function() {
  rnoadm.login.dummy_.value = rnoadm.login.password_.value;
};


/**
 */
rnoadm.login.form_.onsubmit = function() {
  rnoadm.login.username = rnoadm.login.username_.value;
  rnoadm.login.password = rnoadm.login.password_.value;
  if (!rnoadm.login.username.length) {
    rnoadm.login.username = rnoadm.login.password = '';
    rnoadm.login.username_.focus();
    return;
  }
  if (rnoadm.login.password.length <= 2) {
    rnoadm.login.username = rnoadm.login.password = '';
    rnoadm.login.password_.focus();
    return;
  }
  window.sessionStorage['rnoadm_username'] = rnoadm.login.username;
  window.sessionStorage['rnoadm_password'] = rnoadm.login.password;
  rnoadm.net.send({'Auth': {
    'U': rnoadm.login.username,
    'P': rnoadm.login.password
  }});
  rnoadm.login.onlogin_();
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
