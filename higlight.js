function highlightJBGX() {
  const regex = /\bJBGX\d+\b/g;

  const elements = document.querySelectorAll("body *:not(script):not(style)");

  elements.forEach(el => {
    el.childNodes.forEach(node => {
      if (node.nodeType === 3) { // text node
        const text = node.nodeValue;

        if (regex.test(text)) {
          const frag = document.createDocumentFragment();
          let lastIndex = 0;

          text.replace(regex, (match, index) => {
            // text before match
            frag.appendChild(document.createTextNode(text.slice(lastIndex, index)));

            // highlighted match
            const mark = document.createElement("mark");
            mark.style.background = "yellow";
            mark.textContent = match;
            frag.appendChild(mark);

            lastIndex = index + match.length;
          });

          // remaining text
          frag.appendChild(document.createTextNode(text.slice(lastIndex)));

          node.replaceWith(frag);
        }
      }
    });
  });
}

highlightJBGX();