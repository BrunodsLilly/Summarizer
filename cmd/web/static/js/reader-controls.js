/**
 * Reader Controls - Global reader functionality
 * Handles bionic text, font size adjustments, and other global reader features
 */

// Global reader state
let currentFontSize = 18;

/**
 * Toggle bionic reading mode on/off
 */
function toggleBionic() {
    const content = document.getElementById('reader-content');
    const button = document.getElementById('bionic-toggle');
    
    if (!content || !button) {
        console.warn('Bionic toggle elements not found');
        return;
    }
    
    let bionicEnabled = button.textContent.includes('Disable');
    
    if (!bionicEnabled) {
        applyBionicText(content);
        button.textContent = 'Disable Bionic Reading';
        button.classList.remove('bg-blue-600', 'hover:bg-blue-700');
        button.classList.add('bg-green-600', 'hover:bg-green-700');
        console.log('Bionic reading enabled');
    } else {
        removeBionicText(content);
        button.textContent = 'Enable Bionic Reading';
        button.classList.remove('bg-green-600', 'hover:bg-green-700');
        button.classList.add('bg-blue-600', 'hover:bg-blue-700');
        console.log('Bionic reading disabled');
    }
}

/**
 * Apply bionic text formatting to an element
 * @param {HTMLElement} element - The element to apply bionic text to
 */
function applyBionicText(element) {
    const walker = document.createTreeWalker(
        element,
        NodeFilter.SHOW_TEXT,
        null,
        false
    );
    
    const textNodes = [];
    let node;
    
    // Collect all text nodes, excluding code blocks
    while (node = walker.nextNode()) {
        if (node.parentElement.tagName !== 'CODE' && node.parentElement.tagName !== 'PRE') {
            textNodes.push(node);
        }
    }
    
    // Process each text node
    textNodes.forEach(textNode => {
        const text = textNode.textContent;
        const words = text.split(/(\s+)/);
        const fragment = document.createDocumentFragment();
        
        words.forEach(word => {
            if (word.trim()) {
                const span = document.createElement('span');
                span.className = 'bionic-word';
                
                const cleanWord = word.replace(/[^\w]/g, '');
                if (cleanWord.length > 1) {
                    const bionicLength = Math.max(1, Math.ceil(cleanWord.length * 0.5));
                    const bionicPart = word.substring(0, bionicLength);
                    const normalPart = word.substring(bionicLength);
                    
                    span.innerHTML = `<span class="bionic">${bionicPart}</span><span class="non-bionic">${normalPart}</span>`;
                } else {
                    span.innerHTML = `<span class="non-bionic">${word}</span>`;
                }
                fragment.appendChild(span);
            } else {
                fragment.appendChild(document.createTextNode(word));
            }
        });
        
        textNode.parentNode.replaceChild(fragment, textNode);
    });
}

/**
 * Remove bionic text formatting from an element
 * @param {HTMLElement} element - The element to remove bionic text from
 */
function removeBionicText(element) {
    const bionicWords = element.querySelectorAll('.bionic-word');
    bionicWords.forEach(word => {
        const text = word.textContent;
        word.parentNode.replaceChild(document.createTextNode(text), word);
    });
}

/**
 * Adjust font size of the reader content
 * @param {number} delta - The amount to change font size by (positive or negative)
 */
function adjustFontSize(delta) {
    currentFontSize += delta * 2;
    currentFontSize = Math.max(12, Math.min(28, currentFontSize));
    
    const content = document.getElementById('reader-content');
    if (content) {
        content.style.fontSize = currentFontSize + 'px';
        console.log(`Font size adjusted to ${currentFontSize}px`);
    }
}

// Make functions globally available (scripts load after DOM so this works immediately)
window.toggleBionic = toggleBionic;
window.adjustFontSize = adjustFontSize;

console.log('Reader controls loaded and functions registered globally');