// Package Editable provides an active toolbar for content-editable textareas
// Version 1.0

// TODO - show formatting selected on toolbar when selection changes in our contenteditable
// TODO - perhaps intercept return key to be sure we get a para, and to be sure we insert newline within code sections, not new code tags + br or similar
// TODO - Clean out more HTML cruft from programs like word/textedit

// On document ready, scan for and activate toolbars associated with contenteditable
DOM.Ready(function(){
  // Activate editable content
  Editable.Activate('.content-editable-toolbar');
});

var Editable = (function() {
return {
  // Activate editable elements with selector s
  Activate:function(s) {
    if (!DOM.Exists(s)) {
      return;
    }
  
    DOM.Each('.content-editable-toolbar',function(toolbar){
        // Store associated elements for access later
        toolbar.buttons = toolbar.querySelectorAll('a');
        var dataEditable = toolbar.getAttribute('data-editable');
        toolbar.editable = DOM.First(DOM.Format("#{0}-editable",dataEditable));
        toolbar.textarea = DOM.First(DOM.Format("#{0}-textarea",dataEditable));
        
        if (toolbar.editable === undefined) {
          toolbar.editable = DOM.Nearest(toolbar,'.content-editable')[0];
        }
        if (toolbar.textarea === undefined) {
          toolbar.textarea = DOM.Nearest(toolbar,'.content-textarea')[0];
        }
        
        // Set textarea to hidden initially
        toolbar.textarea.style.display = 'none';
        
        // Listen to a form submit, and call updateContent to make sure 
        // our textarea in the form is up to date with the latest content
        toolbar.textarea.form.addEventListener('submit', function(e) {
            Editable.updateContent(toolbar.editable,toolbar.textarea,false);
            return false;
        });
      
        // Intercept paste on editable and remove complex html before it is pasted in
        toolbar.editable.addEventListener('input', function(e) {
          Editable.cleanHTMLElements(this);
        });
        
        
        // Listen to button clicks within the toolbar
        DOM.ForEach(toolbar.buttons,function(el,i){
          el.addEventListener('click', function(e) {
            var cmd = this.id;
            var insert = "";
             
             switch (cmd){
             case "showCode":
                 Editable.updateContent(toolbar.editable,toolbar.textarea,true);
             break;
             case "createLink": 
                 insert = prompt("Supply the web URL to link to");
                 // Prefix url with http:// if no scheme supplied
                 if (insert.indexOf('http') !== 0) {
                     insert = "http://" + insert;
                 }
             break;
             case "formatblock": 
                 insert = this.getAttribute('data-format');
             break;
             default:
             break;
             }
             
            if (cmd.length > 0) {
                 document.execCommand(cmd,false,insert);
            }
             
            // Find and remove evil html created by browsers 
            var sel = Editable.getSelectionParentElement();
         
            if (sel !== null) {
               // Clean align stuff
               Editable.cleanAlign(cmd,sel);
               Editable.cleanHTMLElements(sel);
               sel.removeAttribute('style');
            } 
            

            return false;
          });
        });
    });
  },// End activate
  
  // cleanAlign
  cleanAlign:function(cmd,el) {
    
    switch (cmd){
     case "justifyCenter": 
         
         if (sel.hasClass('align-center')) {
             sel.removeClass('align-center');
         } else {
             sel.addClass('align-center');
         }
         
         sel.removeClass('align-left').removeClass('align-right');
         sel.removeAttr('style');
     break;
     case "justifyLeft": 
         if (sel.hasClass('align-left')) {
             sel.removeClass('align-left');
         } else {
             sel.addClass('align-left');
         }
     
         sel.removeClass('align-center').removeClass('align-right');
         sel.removeAttribute('style');
         
         
     break;
     case "justifyRight": 

         if (sel.hasClass('align-right')) {
             sel.removeClass('align-right');
         } else {
             sel.addClass('align-right');
         }

         sel.removeClass('align-center').removeClass('align-left');
         sel.removeAttribute('style');
     break;
     }
  },
  
  
  
  // CleanHTML is used to clean the html content from the contenteditable before it is assigned to the textarea
  cleanHTML:function(html) {
    html = html.replace(/<\/?span>/gi,'');// Remove all empty span tags
    html = html.replace(/<\/?font [^>]*>/gi,'');// Remove ALL font tags
    html = html.replace(/&nbsp;</gi,' <');// this is sometimes required, be careful
    html = html.replace(/<p><\/p>/gi,'\n');// pretty format but remove empty paras
    html = html.replace(/<br><\/li>/gi,'<\/li>');
    
    // Remove comments and other MS cruft
    html = html.replace(/<!--[\w\d\[\]\s<\/>:.!="*{};-]*-->/gi,'');
    html = html.replace(/ class\=\"MsoNormal\"/gi,'');
    html = html.replace(/<p><o:p> <\/o:p><\/p>/gi,'');
    
    // Pretty printing elements which follow on from one another
    html = html.replace(/><(li|ul|ol|p|h\d|\/ul|\/ol)>/gi,'>\n<$1>');
    
    return html;
  },
  
  // cleanHTMLElements removes certain attributes which are usually full of junk (style, color etc)
  cleanHTMLElements:function(el) {
    // Browsers tend to use style attributes to add all sorts of awful stuff to the html
    // No inline styles allowed
    DOM.ForEach(el.querySelectorAll('p, div, b, i, h1, h2, h3, h4, h5, h6'),function(e){e.removeAttribute('style');});
    DOM.ForEach(el.querySelectorAll('span'),function(e){e.removeAttribute('style');e.removeAttribute('lang');});
    DOM.ForEach(el.querySelectorAll('font'),function(e){e.removeAttribute('color');});
  },
  
  // If textarea visible update the content editable with new html
  // else update the textarea with new content editable html
  updateContent:function(editable,textarea,toggle) {
      var html = '';
      if (textarea.style.display !== 'none') {
          html = textarea.value;
          editable.innerHTML = html;
          if (toggle){
              editable.style.display = '';
              textarea.style.display = 'none';
          }
      } else {
          html = editable.innerHTML;
          // Cleanup the html by removing plain spans
          html = Editable.cleanHTML(html);
          textarea.value = html;
          if (toggle){
              editable.style.display = 'none';
              textarea.style.display = '';
          }
      }

  },
  
  // The closest parent element which encloses the entire text selection
  getSelectionParentElement:function() {
      var p = null, sel;
      if (window.getSelection) {
          sel = window.getSelection();
          if (sel.rangeCount) {
              p = sel.getRangeAt(0).commonAncestorContainer;
              if (p.nodeType != 1) {
                  p = p.parentNode;
              }
          }
      } else if ( (sel = document.selection) && sel.type != "Control") {
          p = sel.createRange().parentElement();
      }
      return p;
  }
  
  };
}());


