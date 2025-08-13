document.addEventListener('DOMContentLoaded', () => {
    const cart = [];
    const cartItemsContainer = document.getElementById('cartItems');
    const cartTotalElement = document.getElementById('cartTotal');
    const emptyCartMessage = document.getElementById('emptyCartMessage');
    const proceedToOrderBtn = document.getElementById('proceedToOrderBtn');
    const categoryFilter = document.getElementById('categoryFilter');

    categoryFilter.addEventListener('change', (event) => {
        const selectedCategory = event.target.value;
        document.querySelectorAll('.category-section').forEach(section => {
            if (selectedCategory === 'all' || section.dataset.category === selectedCategory) {
                section.style.display = 'block';
            } else {
                section.style.display = 'none';
            }
        });
    });

    document.querySelectorAll('.add-to-cart-btn').forEach(button => {
        button.addEventListener('click', (event) => {
            const itemId = event.target.dataset.itemId;
            const itemName = event.target.dataset.itemName;
            const itemPrice = parseFloat(event.target.dataset.itemPrice);
            const quantityInput = document.getElementById(`quantity-${itemId}`);
            const instructionsInput = document.getElementById(`instructions-${itemId}`);
            const quantity = parseInt(quantityInput.value);
            const instructions = instructionsInput.value.trim();

            if (quantity <= 0) {
                alert('Please enter a valid quantity.');
                return;
            }

            const existingItem = cart.find(item => item.id === itemId);

            if (existingItem) {
                existingItem.quantity += quantity;
                existingItem.instructions = instructions;
            } else {
                cart.push({ id: itemId, name: itemName, price: itemPrice, quantity: quantity, instructions: instructions });
            }

            updateCartDisplay();
            quantityInput.value = 1;
            instructionsInput.value = '';
        });
    });

    function updateCartDisplay() {
        cartItemsContainer.innerHTML = '';
        let total = 0;

        if (cart.length === 0) {
            emptyCartMessage.style.display = 'block';
            proceedToOrderBtn.disabled = true;
        } else {
            emptyCartMessage.style.display = 'none';
            proceedToOrderBtn.disabled = false;
            cart.forEach((item, index) => {
                total += item.quantity * item.price;
                const cartItemDiv = document.createElement('div');
                cartItemDiv.className = 'cart-item';
                cartItemDiv.innerHTML = `
                    <div>
                        <h5 class="cart-item-name">${item.name} (${item.quantity}x)</h5>
                        <p class="cart-item-price">$${item.price.toFixed(2)} each</p>
                        ${item.instructions ? `<p class="cart-item-instructions">Instructions: ${item.instructions}</p>` : ''}
                    </div>
                    <button class="remove-item-btn action-button" data-index="${index}">&times; Remove</button>
                `;
                cartItemsContainer.appendChild(cartItemDiv);
            });
        }
        cartTotalElement.textContent = `$${total.toFixed(2)}`;
    }

    cartItemsContainer.addEventListener('click', (event) => {
        if (event.target.classList.contains('remove-item-btn')) {
            const indexToRemove = parseInt(event.target.dataset.index);
            cart.splice(indexToRemove, 1);
            updateCartDisplay();
        }
    });

    const orderForm = document.getElementById('orderForm');
    orderForm.addEventListener('submit', (event) => {
        event.preventDefault();

        if (cart.length === 0) {
            alert('Your cart is empty. Please add items before proceeding to order.');
            return;
        }

        const createHiddenInput = (name, value) => {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = name;
            input.value = value;
            return input;
        };

        cart.forEach(item => {
            orderForm.appendChild(createHiddenInput('itemId', item.id));
            orderForm.appendChild(createHiddenInput('quantity', item.quantity));
            orderForm.appendChild(createHiddenInput('instruction', item.instructions));
        });

        orderForm.submit();
    });

    updateCartDisplay();
});