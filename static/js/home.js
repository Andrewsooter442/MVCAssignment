document.addEventListener('DOMContentLoaded', () => {
    const cart = [];
    const cartItemsContainer = document.getElementById('cartItems');
    const cartTotalElement = document.getElementById('cartTotal');
    const emptyCartMessage = document.getElementById('emptyCartMessage');
    const proceedToOrderBtn = document.getElementById('proceedToOrderBtn');
    const categoryFilter = document.getElementById('categoryFilter');
    const orderForm = document.getElementById('orderForm');

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

    document.querySelectorAll('.quantity-btn').forEach(button => {
        button.addEventListener('click', (event) => {
            const clickedButton = event.currentTarget;
            const itemId = clickedButton.dataset.itemId;
            const input = document.getElementById(`quantity-${itemId}`);

            if (!input) return;

            let currentValue = parseInt(input.value, 10);
            const minValue = parseInt(input.min, 10);

            if (clickedButton.classList.contains('quantity-plus')) {
                currentValue++;
            } else if (clickedButton.classList.contains('quantity-minus')) {
                currentValue--;
            }

            if (currentValue < minValue) {
                currentValue = minValue;
            }

            input.value = currentValue;
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

            const existingItem = cart.find(item => item.id === itemId && item.instructions === instructions);

            if (existingItem) {
                existingItem.quantity += quantity;
            } else {
                cart.push({ id: itemId, name: itemName, price: itemPrice, quantity: quantity, instructions: instructions });
            }

            updateCartDisplay();
            quantityInput.value = 1;
            if (instructionsInput) instructionsInput.value = '';
        });
    });

    function updateCartDisplay() {
        cartItemsContainer.innerHTML = '';
        let total = 0;

        if (cart.length === 0) {
            if(emptyCartMessage) emptyCartMessage.style.display = 'block';
            if(proceedToOrderBtn) proceedToOrderBtn.disabled = true;
        } else {
            if(emptyCartMessage) emptyCartMessage.style.display = 'none';
            if(proceedToOrderBtn) proceedToOrderBtn.disabled = false;

            cart.forEach((item, index) => {
                total += item.quantity * item.price;
                const cartItemDiv = document.createElement('div');
                cartItemDiv.className = 'cart-item';
                cartItemDiv.innerHTML = `
                    <div>
                        <h5 class="cart-item-name">${item.name} (${item.quantity}x)</h5>
                        <p class="cart-item-price">$${item.price.toFixed(2)} each</p>
                        ${item.instructions ? `<p class="cart-item-instructions"><strong>Instructions:</strong> ${item.instructions}</p>` : ''}
                    </div>
                    <button class="remove-item-btn action-button" data-index="${index}">&times;</button>
                `;
                cartItemsContainer.appendChild(cartItemDiv);
            });
        }
        if(cartTotalElement) cartTotalElement.textContent = `$${total.toFixed(2)}`;
    }

    cartItemsContainer.addEventListener('click', (event) => {
        if (event.target.classList.contains('remove-item-btn')) {
            const indexToRemove = parseInt(event.target.dataset.index);
            cart.splice(indexToRemove, 1);
            updateCartDisplay();
        }
    });

    if (orderForm) {
        orderForm.addEventListener('submit', (event) => {
            event.preventDefault();

            if (cart.length === 0) {
                alert('Your cart is empty. Please add items before placing an order.');
                return;
            }

            orderForm.querySelectorAll('input[type="hidden"]').forEach(input => {
                if (input.name !== 'tableNumber') {
                    input.remove();
                }
            });

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
    }

    updateCartDisplay();
});