// Creates a dual-slider. See Readme for usage
class DualSlider {
    constructor(containerId, options = {}) {
        const {
            minValue = 0, maxValue = 100, initialLeft = minValue, initialRight = maxValue
        } = options;

        this.containerId = containerId;
        this.minValue = minValue;
        this.maxValue = maxValue;
        this.container = document.getElementById(containerId);
        this.activeSlider = null;

        this.init();
        this.setInitialValues(initialLeft, initialRight);
    }

    init() {
        // Create HTML structure
        this.container.innerHTML = `
                    <div class="display-container">
                        <div class="value-display left-value">0</div>
                        <div class="value-display right-value">100</div>
                    </div>
                    <div class="slider-container">
                        <div class="slider-track"></div>
                        <div class="slider-range"></div>
                        <div class="slider left-slider"></div>
                        <div class="slider right-slider"></div>
                    </div>
                `;

        // Get DOM elements
        this.leftSlider = this.container.querySelector('.left-slider');
        this.rightSlider = this.container.querySelector('.right-slider');
        this.sliderRange = this.container.querySelector('.slider-range');
        this.leftValueDisplay = this.container.querySelector('.left-value');
        this.rightValueDisplay = this.container.querySelector('.right-value');

        // Bind event listeners
        this.leftSlider.addEventListener('mousedown', this.handleMouseDown.bind(this));
        this.rightSlider.addEventListener('mousedown', this.handleMouseDown.bind(this));
    }

    setInitialValues(leftValue, rightValue) {
        const leftPos = (leftValue - this.minValue) / (this.maxValue - this.minValue) * 280;
        const rightPos = (rightValue - this.minValue) / (this.maxValue - this.minValue) * 280;

        this.leftSlider.style.left = leftPos + 'px';
        this.rightSlider.style.left = rightPos + 'px';
        this.updateSliderRange();
    }

    handleMouseDown(e) {
        this.activeSlider = e.target;
        document.addEventListener('mousemove', this.handleMouseMove.bind(this));
        document.addEventListener('mouseup', this.handleMouseUp.bind(this));
    }

    handleMouseMove(e) {
        if (!this.activeSlider) return;

        const containerRect = this.container.querySelector('.slider-container').getBoundingClientRect();
        let newLeft = e.clientX - containerRect.left - 10;
        newLeft = Math.max(0, Math.min(newLeft, 280));

        if (this.activeSlider.classList.contains('left-slider')) {
            const rightPos = parseInt(this.rightSlider.style.left);
            newLeft = Math.min(newLeft, rightPos);
            this.leftSlider.style.left = newLeft + 'px';
        } else {
            const leftPos = parseInt(this.leftSlider.style.left);
            newLeft = Math.max(newLeft, leftPos);
            this.rightSlider.style.left = newLeft + 'px';
        }

        this.updateSliderRange();
    }

    handleMouseUp() {
        this.activeSlider = null;
        document.removeEventListener('mousemove', this.handleMouseMove.bind(this));
        document.removeEventListener('mouseup', this.handleMouseUp.bind(this));
    }

    updateSliderRange() {
        const leftPos = parseInt(this.leftSlider.style.left || '0');
        const rightPos = parseInt(this.rightSlider.style.left || '280');

        this.sliderRange.style.left = leftPos + 10 + 'px';
        this.sliderRange.style.width = (rightPos - leftPos) + 'px';

        const leftValue = Math.round((leftPos / 280) * (this.maxValue - this.minValue) + this.minValue);
        const rightValue = Math.round((rightPos / 280) * (this.maxValue - this.minValue) + this.minValue);

        this.leftValueDisplay.textContent = leftValue;
        this.rightValueDisplay.textContent = rightValue;

        // Dispatch custom event with slider ID
        const event = new CustomEvent('sliderChange', {
            detail: {
                sliderId: this.containerId, leftValue, rightValue
            }
        });
        document.dispatchEvent(event);
    }
}
