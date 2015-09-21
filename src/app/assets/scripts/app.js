DOM.Ready(function(){
  // Show/Hide elements with selector in attribute data-show
  ActivateShowlinks();
  ActivatePackagelinks();
  // Perform AJAX post on click on method=post|delete anchors
  ActivateMethodLinks();
  // Submit forms of class .filter-form when filter fields change
  ActivateFilterFields();
  // Activate pager links (on home)
  ActivatePagerlinks();
  // Activate the map overlay
  ActivateMap();
  
  // Activate highlight.js on all pages
  hljs.initHighlightingOnLoad();
});

// Perform AJAX post on click on method=post|delete anchors
function ActivateMethodLinks() {
  DOM.On('a[method="post"], a[method="delete"]','click',function(e){
    // Confirm action before delete
    if (this.getAttribute('method') == 'delete') {
      if (!confirm('Are you sure you want to delete this item, this action cannot be undone?')) { 
          return false;
      }
    }

    // Perform a post to the specified url (href of link)
    var url       = this.getAttribute('href');
    var redirect  = this.getAttribute('data-redirect');
    
    DOM.Post(url,function(){
      if (redirect.length > 0) {
        // If we have a redirect, redirect to it after the link is clicked
        window.location = redirect;
      } else {
        // If no redirect supplied, we just reload the current screen
        window.location.reload();
      }
    },function(){
      console.log("#error POST to"+url+"failed");
    });
    
    
    return false;
  });
}

// Submit forms of class .filter-form when filter fields change
function ActivateFilterFields() {
  DOM.On('.filter-form .field select, .filter-form .field input','change',function(e){
      this.form.submit();
  });
}

// Show/Hide elements with selector in attribute href - do this with a hidden class name
function ActivateShowlinks() {
  DOM.On('.show','click',function(e){
    var selector = this.getAttribute('href');
      DOM.Each(selector,function(el,i){
        if (el.className != 'hidden') {
           el.className = 'hidden';
        } else {
           el.className = el.className.replace(/hidden/gi,'');
        }
      });
      
      return false;
  });
}

// Show/Hide elements with selector in attribute href - do this with a hidden class name
function ActivatePackagelinks() {
  DOM.On('.feature_dot','click',function(e){
    var selector = this.getAttribute('href');
    
     // First hide all packages
     DOM.Each(".package",function(el,i){
        el.className = 'package hidden';
     });
     
     DOM.First(selector).className = 'package';
     
     return false;
  });
}


// Show/Hide elements with selector in attribute data-show
function ActivatePagerlinks() {
  DOM.On('.pager li a','click',function(e){
    // Update pager selection stage
    DOM.ForEach(this.parentNode.parentNode.querySelectorAll('a'),function(el){
      el.className = '';
    });
    this.className = 'selected';
    
    // Update pager elements - first find the element linked (href=#id)
    var el = DOM.First(this.getAttribute('href'));
    // Hide everything with the same first class
    DOM.Hide("."+el.className.split(" ")[0]);
    // Show the element concerned
    el.style.display = 'block';
    
    e.preventDefault();
    return false;
  });
}


function ActivateMap() {
  DOM.On('.contact_iframe_container','click',function(e){
    this.style.pointerEvents = 'none';/* ignore this after first click on it */
  });
}
