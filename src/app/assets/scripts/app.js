DOM.Ready(function() {

    // Perform AJAX post on click on method=post|delete anchors
    ActivateMethodLinks();

    // Show/Hide elements with selector in attribute data-show
    ActivateShowlinks();

    // Submit forms of class .filter-form when filter fields change
    ActivateFilterFields();

    // Insert CSRF tokens into forms
    ActivateForms();

});

// Perform AJAX post on click on method=post|delete anchors
function ActivateMethodLinks() {
    DOM.On('a[method="post"], a[method="delete"]', 'click', function(e) {
        // Confirm action before delete
        if (this.getAttribute('method') == 'delete') {
            if (!confirm('Are you sure you want to delete this item, this action cannot be undone?')) {
                return false;
            }
        }

        // Get authenticity token from head of page
        var token = authenticityToken();

        // Perform a post to the specified url (href of link)
        var url = this.getAttribute('href');
        var data = "authenticity_token=" + token

        DOM.Post(url, data, function(request) {
            // Use the response url to redirect 
            window.location = request.responseURL
        }, function(request) {
            console.log("error with request:", request)
                // Use the response url to redirect even if not found
            window.location = request.responseURL
        });

        e.preventDefault();
        return false;
    });


    DOM.On('a[method="back"]', 'click', function(e) {
        history.back(); // go back one step in history
        e.preventDefault();
        return false;
    });

}


// Insert an input into every form with js to include the csrf token.
// this saves us having to insert tokens into every form.
function ActivateForms() {
    // Get authenticity token from head of page
    var token = authenticityToken();

    DOM.Each('form', function(f) {

        // Create an input element 
        var csrf = document.createElement("input");
        csrf.setAttribute("name", "authenticity_token");
        csrf.setAttribute("value", token);
        csrf.setAttribute("type", "hidden");

        //Append the input
        f.appendChild(csrf);
    });
}

// Submit forms of class .filter-form when filter fields change
function ActivateFilterFields() {
    DOM.On('.filter-form .field select, .filter-form .field input', 'change', function(e) {
        this.form.submit();
    });
}

// Show/Hide elements with selector in attribute href - do this with a hidden class name
function ActivateShowlinks() {
    DOM.On('.show', 'click', function(e) {
        var selector = this.getAttribute('href');
        DOM.Each(selector, function(el, i) {
            if (el.className != 'hidden') {
                el.className = 'hidden';
            } else {
                el.className = el.className.replace(/hidden/gi, '');
            }
        });

        return false;
    });
}

function authenticityToken() {
    // Collect the authenticity token from meta tags in header
    var meta = DOM.First("meta[name='authenticity_token']")
    if (meta === undefined) {
        e.preventDefault();
        return ""
    }
    return meta.getAttribute('content');
}