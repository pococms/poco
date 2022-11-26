  // Search for "USER JAVASCRIPT CODE" to see where user Javascript should go.
  // Document.ready working back to IE 8.
  // From https://stackoverflow.com/questions/2304941/what-is-the-non-jquery-equivalent-of-document-ready
  function ready(fn) {
    if (typeof fn !== 'function') {
      throw new Error('Argument passed to ready should be a function');
    }

    if (document.readyState != 'loading') {
      fn();
    } else if (document.addEventListener) {
      document.addEventListener('DOMContentLoaded', fn, {
        once: true // A boolean value indicating that the listener should be invoked at most once after being added. If true, the listener would be automatically removed when invoked.
      });
    } else {
      document.attachEvent('onreadystatechange', function() {
        if (document.readyState != 'loading')
          fn();
      });
    }
  }

  // Demonstrate error checking 
  //try {
  //  window.ready(5);
  //} catch (ex) {
  //  console.log(ex.message);
  //}


  // HTML Document has been loaded. This is where it's safe 
  // for you to execute code. For modularity's sake this 
  // version calls an arbitrary function.
  window.ready(function() {
  // startHere() is just an arbitrary name. Can be anything.
  startHere();
  })

  // --- START USER JAVASCRIPT CODE
  function startHere() {
  }
  // --- END USER JAVASCRIPT CODE
;


