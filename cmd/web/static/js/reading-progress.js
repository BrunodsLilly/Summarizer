/**
 * Reading Progress Tracker
 * Handles word counting, reading time estimation, and scroll-based progress tracking
 */

/**
 * Initialize reading progress tracking for a summary
 * This function should be called when new content is loaded
 */
function initializeReadingProgress() {
    console.log('Initializing reading progress for new summary...');
    
    let totalWords = 0;
    const readingSpeed = 200; // words per minute (average reading speed)
    
    const content = document.getElementById('reader-content');
    if (!content) {
        console.log('Reader content not found');
        return;
    }
    
    // Calculate word count
    const text = content.textContent || content.innerText;
    const words = text.trim().split(/\s+/).filter(word => word.length > 0);
    totalWords = words.length;
    
    console.log('Reading progress initialized:', {
        contentFound: !!content,
        textLength: text.length,
        totalWords: totalWords
    });
    
    if (totalWords === 0) {
        console.log('No words found in content');
        return;
    }
    
    // Calculate reading time
    const totalMinutes = Math.max(1, Math.ceil(totalWords / readingSpeed));
    
    // Update UI elements
    const wordCountEl = document.getElementById('word-count');
    const readingTimeEl = document.getElementById('reading-time');
    const timeRemainingEl = document.getElementById('time-remaining');
    
    if (wordCountEl) wordCountEl.textContent = `${totalWords.toLocaleString()} words`;
    if (readingTimeEl) readingTimeEl.textContent = `~${totalMinutes} min read`;
    if (timeRemainingEl) timeRemainingEl.textContent = `~${totalMinutes} min remaining`;
    
    // Set initial font size
    content.style.fontSize = (window.currentFontSize || 18) + 'px';
    
    // Setup scroll tracking
    setupScrollTracking(totalWords, readingSpeed);
}

/**
 * Setup scroll-based progress tracking
 * @param {number} totalWords - Total word count in the content
 * @param {number} readingSpeed - Reading speed in words per minute
 */
function setupScrollTracking(totalWords, readingSpeed) {
    const content = document.getElementById('reader-content');
    const progressFill = document.getElementById('progress-fill');
    const progressPercent = document.getElementById('progress-percent');
    const timeRemaining = document.getElementById('time-remaining');
    
    if (!content || !progressFill) {
        console.log('Missing elements for scroll tracking');
        return;
    }
    
    /**
     * Update progress based on current scroll position
     */
    function updateProgress() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;
        
        // Simple document-based progress calculation
        const maxScroll = Math.max(0, documentHeight - windowHeight);
        let progress = 0;
        
        if (maxScroll > 0) {
            progress = Math.min(100, (scrollTop / maxScroll) * 100);
        } else {
            // Document fits in viewport
            progress = 100;
        }
        
        // Update UI elements
        progressFill.style.width = `${progress}%`;
        if (progressPercent) progressPercent.textContent = `${Math.round(progress)}% complete`;
        
        // Update time remaining
        const remainingWords = Math.ceil(totalWords * (100 - progress) / 100);
        const remainingMinutes = Math.max(0, Math.ceil(remainingWords / readingSpeed));
        if (timeRemaining) {
            timeRemaining.textContent = remainingMinutes > 0 ? 
                `~${remainingMinutes} min remaining` : 'Complete!';
        }
    }
    
    // Throttled scroll listener for performance
    let ticking = false;
    function onScroll() {
        if (!ticking) {
            requestAnimationFrame(() => {
                updateProgress();
                ticking = false;
            });
            ticking = true;
        }
    }
    
    // Add event listeners
    window.addEventListener('scroll', onScroll);
    window.addEventListener('resize', updateProgress);
    
    // Initial calculation
    updateProgress();
    
    console.log('Scroll tracking setup complete');
}

// Make function globally available for template usage
window.initializeReadingProgress = initializeReadingProgress;

// Auto-initialize if reader content is already present
document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('reader-content')) {
        console.log('Auto-initializing reading progress');
        initializeReadingProgress();
    }
});

// Also initialize if the script loads after DOM is ready
if (document.readyState === 'loading') {
    // DOM hasn't finished loading yet
    document.addEventListener('DOMContentLoaded', function() {
        if (document.getElementById('reader-content')) {
            console.log('Auto-initializing reading progress (late load)');
            initializeReadingProgress();
        }
    });
} else {
    // DOM is already ready
    if (document.getElementById('reader-content')) {
        console.log('Auto-initializing reading progress (immediate)');
        initializeReadingProgress();
    }
}

console.log('Reading progress module loaded');