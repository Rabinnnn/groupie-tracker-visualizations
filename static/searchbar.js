document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("search-input");
    const suggestionsList = document.getElementById("suggestions");
    const searchButton = document.getElementById("search-button");
    let cachedSuggestions = [];
    let currentFocus = -1;
  
    // Fetch and cache suggestions on page load
    fetch('/search-suggestions?q=')
      .then(response => response.json())
      .then(suggestions => {
        cachedSuggestions = suggestions;
      })
      .catch(err => {
        console.error("Error fetching initial suggestions:", err);
      });
  
    // Debounce function
    function debounce(func, delay) {
      let timeoutId;
      return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(this, args), delay);
      };
    }
  
    // Filter suggestions based on input
    function filterSuggestions(query) {
      query = query.toLowerCase();
      return cachedSuggestions.filter(suggestion => 
        suggestion.toLowerCase().includes(query)
      );
    }
  
    // Display suggestions
    function displaySuggestions(suggestions) {
      suggestionsList.innerHTML = "";
      const seenSuggestions = new Set();
      currentFocus = -1;
      
      suggestions.forEach(suggestion => {
        const nameTypeCombo = suggestion.toLowerCase();
        
        if (!seenSuggestions.has(nameTypeCombo)) {
          seenSuggestions.add(nameTypeCombo);
  
          const li = document.createElement("li");
          li.textContent = suggestion;
          li.addEventListener("click", () => {
            searchInput.value = suggestion.split(" - ")[0];
            performSearch(searchInput.value);
          });
          li.addEventListener("mouseover", () => {
            removeActive(suggestionsList.getElementsByTagName("li"));
            li.classList.add("active");
          });
          suggestionsList.appendChild(li);
        }
      });
    }
  
    // Throttled search function
    const throttledSearch = debounce((query) => {
      if (query === "") {
        suggestionsList.innerHTML = "";
        return;
      }
  
      const filteredSuggestions = filterSuggestions(query);
      displaySuggestions(filteredSuggestions);
    }, 100);
  
    function performSearch(query) {
      window.location.href = `/?query=${encodeURIComponent(query)}`;
    }
  
    function addActive(x) {
      if (!x) return false;
      removeActive(x);
      if (currentFocus >= x.length) currentFocus = 0;
      if (currentFocus < 0) currentFocus = (x.length - 1);
      x[currentFocus].classList.add("active");
      // Auto-scroll
      x[currentFocus].scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
        inline: 'start'
      });
    }
  
    function removeActive(x) {
      for (let i = 0; i < x.length; i++) {
        x[i].classList.remove("active");
      }
    }
  
    searchInput.addEventListener("input", function() {
      const query = searchInput.value.trim();
      throttledSearch(query);
    });
  
    searchInput.addEventListener("keydown", function(event) {
      const x = suggestionsList.getElementsByTagName("li");
      if (event.key === "ArrowDown") {
        currentFocus++;
        addActive(x);
      } else if (event.key === "ArrowUp") {
        currentFocus--;
        addActive(x);
      } else if (event.key === "Enter") {
        event.preventDefault();
        if (currentFocus > -1) {
          if (x) x[currentFocus].click();
        } else {
          performSearch(searchInput.value.trim());
        }
      }
    });
  
    searchButton.addEventListener("click", function() {
      performSearch(searchInput.value.trim());
    });
  
    document.addEventListener("click", function(event) {
      if (!searchInput.contains(event.target) && !suggestionsList.contains(event.target)) {
        suggestionsList.innerHTML = "";
      }
    });
  });