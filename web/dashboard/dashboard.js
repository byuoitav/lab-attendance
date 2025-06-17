function createSquare(positionClass, onTap) {
    const square = document.createElement('div');
    square.className = `square ${positionClass}`;

    // Show the square on touch or mouse down
    const show = () => square.style.opacity = 0.2;
    const hide = () => square.style.opacity = 0;

    // Event handling
    square.addEventListener('mousedown', e => { show(); onTap(); });
    square.addEventListener('mouseup', hide);
    square.addEventListener('touchstart', e => { show(); onTap(); });
    square.addEventListener('touchend', hide);

    document.body.appendChild(square);
    return square;
}

let topLeft, topRight, bottomLeft, bottomRight;

topLeft = createSquare('top-left', () => {
    if (!topRight) {
        topRight = createSquare('top-right', () => {
            if (!bottomLeft) {
                bottomLeft = createSquare('bottom-left', () => {
                    if (!bottomRight) {
                        bottomRight = createSquare('bottom-right', () => {
                            window.location.href = 'http://localhost:10000/dashboard/overview';
                        });
                    }
                });
            }
        });
    }
});

// Clean up all but topLeft after 20 seconds
setTimeout(() => {
    [topRight, bottomLeft, bottomRight].forEach(sq => {
        if (sq && sq.parentElement) sq.remove();
    });
    topRight = bottomLeft = bottomRight = null;
}, 20000);