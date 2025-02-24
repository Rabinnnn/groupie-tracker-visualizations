#### Dual Slider

1. Create a container div with a unique ID:

    ```html
    <div id="mySlider1" class="dual-slider-container"></div>
    ```

2. Create a new instance with options:

    ```javascript 
    const slider = new DualSlider('mySlider1', {
        minValue: 0,
        maxValue: 100,
        initialLeft: 20,
        initialRight: 80
    });
    ```

3. Listen for changes:
    
    ```js
    document.addEventListener('sliderChange', (e) => {
        const { sliderId, leftValue, rightValue } = e.detail;
        // Now you can handle each slider separately based on sliderId
        if (sliderId === 'mySlider1') {
            console.log('Slider 1:', leftValue, rightValue);
        }
    });
    ```      
   
#### Create Multiple instances

```js

const slider1 = new DualSlider('slider1', {
    minValue: 0, maxValue: 100, initialLeft: 20, initialRight: 80
});

const slider2 = new DualSlider('slider2', {
    minValue: 0, maxValue: 200, initialLeft: 50, initialRight: 150
});

// Example of listening to slider changes
document.addEventListener('sliderChange', (e) => {
    console.log('Slider:', e.detail.sliderId, 'Left:', e.detail.leftValue, 'Right:', e.detail.rightValue);
});

```