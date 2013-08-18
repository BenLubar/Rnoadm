goog.provide('rnoadm.login');

goog.require('rnoadm.net');


/**
 * @type {string}
 */
rnoadm.login.username = sessionStorage['rnoadm_username'] || '';


/**
 * @type {string}
 */
rnoadm.login.password = sessionStorage['rnoadm_password'] || '';


/**
 * @type {HTMLFormElement}
 * @private
 * @const
 */
rnoadm.login.form_ = document.querySelector('form');


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.username_ = rnoadm.login.form_.querySelector('#username');


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.password_ = rnoadm.login.form_.querySelector('#password');


/**
 * @type {HTMLInputElement}
 * @private
 * @const
 */
rnoadm.login.dummy_ = rnoadm.login.form_.querySelector('#password2');


/**
 * @param {...string} var_args
 * @private
 */
rnoadm.net.admin_ = function(var_args) {
  rnoadm.net.send({'Admin': [].slice.call(arguments)});
};

goog.exportSymbol('admin', rnoadm.net.admin_);


/**
 * @private
 */
rnoadm.login.onlogin_ = function() {
  var parent = rnoadm.login.form_.parentNode;
  parent.removeChild(rnoadm.login.form_);
  parent.style.overflow = 'hidden';
  parent.style.fontSize = '0';
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
  sessionStorage['rnoadm_username'] = rnoadm.login.username;
  sessionStorage['rnoadm_password'] = rnoadm.login.password;
  rnoadm.net.send({'Auth': {
    'U': rnoadm.login.username,
    'P': rnoadm.login.password
  }});
  rnoadm.login.onlogin_();
};

// vim: set tabstop=2 shiftwidth=2 expandtab:
