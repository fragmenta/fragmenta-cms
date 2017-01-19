// Package DOM provides functions to replace the use of jquery in 1.4KB of js
// See http://youmightnotneedjquery.com/ for more if required
var DOM = (function() {
    return {
        // Apply a function on document ready
        Ready: function(f) {
            if (document.readyState != 'loading') {
                f();
            } else {
                document.addEventListener('DOMContentLoaded', f);
            }
        },

        // Return a NodeList of nearest elements matching selector, 
        // checking children, siblings or parents of el
        Nearest: function(el, s) {

            // Start with this element, then walk up the tree till 
            // we find a child which matches selector or we run out of elements
            while (el !== undefined && el !== null) {
                var nearest = el.querySelectorAll(s);
                if (nearest.length > 0) {
                    return nearest;
                }
                el = el.parentNode;
            }

            return []; // return empty array
        },

        // FIXME - perhaps adjust all to operate on either selector or an element?

        // Attribute returns either an attribute value or an empty string (if null)
        Attribute: function(el, a) {
            if (el.getAttribute(a) === null) {
                return ''
            }
            return el.getAttribute(a)
        },

        // HasClass returns true if this element has this className
        HasClass: function(el, c) {
            var regexp = new RegExp("\\b" + c + "\\b", 'gi');
            return regexp.test(el.className);
        },

        // AddClass Adds the given className from el.className
        AddClass: function(s, c) {
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    if (!DOM.HasClass(el, c)) {
                        el.className = el.className + ' ' + c;
                    }
                });
            } else {
                if (!DOM.HasClass(s, c)) {
                    s.className = s.className + ' ' + c;
                }
            }
        },

        // RemoveClass removes the given className from el.className
        RemoveClass: function(s, c) {
            var regexp = new RegExp("\\b" + c + "\\b", 'gi');
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    el.className = el.className.replace(regexp, '')
                });
            } else {
                s.className = s.className.replace(regexp, '')
            }
        },

        // Format returns the format string with the indexed arguments substituted
        // Formats are of the form - "{0} {1}" which uses variables 0 and 1 respectively
        Format: function(f) {
            for (var i = 1; i < arguments.length; i++) {
                var regexp = new RegExp('\\{' + (i - 1) + '\\}', 'gi');
                f = f.replace(regexp, arguments[i]);
            }
            return f;
        },


        // Apply a function to elements of an array
        ForEach: function(a, f) {
            Array.prototype.forEach.call(a, f);
        },


        // Return true if any element match selector
        Exists: function(s) {
            return (document.querySelector(s) !== null);
        },

        // Return a NodeList of elements matching selector
        All: function(s) {
            return document.querySelectorAll(s);
        },


        // Return the first in the NodeList of elements matching selector - may return nil
        First: function(s) {
            return DOM.All(s)[0];
        },

        // Apply a function to elements matching selector, return true to break
        Each: function(s, f) {
            var a = DOM.All(s);
            for (i = 0; i < a.length; ++i) {
              f(a[i], i);
            }
        },


        // Hidden returns true if this element is hidden
        Hidden: function(s) {
            if (typeof s === "string") {
                return (DOM.First(s).style.display == 'none');
            } else {
                return s.style.display == 'none';
            }

        },

        // Hide elements matching selector 
        Hide: function(s) {
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    el.style.display = 'none';
                });
            } else {
                s.style.display = 'none';
            }
        },

        // Show elements matching selector
        Show: function(s) {
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    el.style.display = '';
                });
            } else {
                s.style.display = '';
            }
        },

        // Toggle the Shown or Hidden value of elements matching selector
        ShowHide: function(s) {
            if (typeof s === "string") {
                DOM.Each(s, function(el, i) {
                    if (el.style.display != 'none') {
                        el.style.display = 'none';
                    } else {
                        el.style.display = '';
                    }
                });
            } else {
                if (s.style.display != 'none') {
                    s.style.display = 'none';
                } else {
                    s.style.display = '';
                }
            }
        },

        // Attach event handlers to all matches for a selector 
        On: function(s, b, f) {
            DOM.Each(s, function(el, i) {
                el.addEventListener(b, f);
            });
        },


        // Ajax - Send to url u the data d, call fs for success, ff for failures
        Post: function(u, d, fs, fe) {
            var request = new XMLHttpRequest();
            request.open('POST', u, true);
            request.onerror = fe;
            request.onload = function() {
                if (request.status >= 200 && request.status < 400) {
                    fs(request);
                } else {
                    fe(request);
                }
            };
            request.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded; charset=UTF-8');
            request.send(d);
        },

        // Ajax - Get the data from url u, call fs for success, ff for failures
        Get: function(u, fs, fe) {
            var request = new XMLHttpRequest();
            request.open('GET', u, true);
            request.onload = function() {
                if (request.status >= 200 && request.status < 400) {
                    fs(request);
                } else {
                    fe();
                }
            };
            request.onerror = fe;
            request.send();
        }

    };

}());