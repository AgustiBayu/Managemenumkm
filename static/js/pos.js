document.addEventListener('DOMContentLoaded', () => {
    // --- DOM Elements ---
    const productGrid = document.getElementById('product-grid');
    const cartTableBody = document.getElementById('cart-items-table-body');
    const productSearch = document.getElementById('product-search');
    const checkoutButton = document.querySelector('.btn-payment');
    const cancelCartButton = document.querySelector('.btn-cancel');
    const clockElement = document.getElementById('header-time');
    const categoryFilterBar = document.getElementById('category-filter-bar');


    // Summary elements
    const summaryTotalItems = document.getElementById('summary-total-items');
    const summaryTotal = document.getElementById('summary-total');
    const summaryTotalPayable = document.getElementById('summary-total-payable');
    // Templates
    const productCardTemplate = document.getElementById('product-card-template');
    const cartItemRowTemplate = document.getElementById('cart-item-row-template');
    // --- State ---
    let allProducts = [];
    let cart = []; // Array of { id, name, code, price, quantity }
    let categories = [];
    let activeCategoryId = 'all';


    // --- Utility Functions ---
    const formatCurrency = (amount) => 
        new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);
    const updateClock = () => {
        if (clockElement) {
            clockElement.textContent = new Date().toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'medium' });
        }
    };
    // --- Render Functions ---
    const renderProducts = (productsToRender) => {
        productGrid.innerHTML = '';
        productsToRender.forEach(product => {
            const card = productCardTemplate.content.cloneNode(true);
            const cardElement = card.querySelector('.product-card');
            cardElement.dataset.productId = product.ID;
            card.querySelector('.product-image').src = product.image_url || '/static/image/placeholder.svg';
            card.querySelector('.product-name').textContent = product.name;
            card.querySelector('.product-category').textContent = product.category ? product.category.category : 'Uncategorized';
            card.querySelector('.product-price').textContent = formatCurrency(product.price);
            card.querySelector('.product-stock').textContent = `Stock: ${product.total_stock}`;
            productGrid.appendChild(card);
        });
    };

    const renderCategories = () => {
        categoryFilterBar.innerHTML = '';
        
        // "All" button
        const allButton = document.createElement('button');
        allButton.textContent = 'All';
        allButton.dataset.categoryId = 'all';
        allButton.classList.toggle('active', activeCategoryId === 'all');
        categoryFilterBar.appendChild(allButton);

        categories.forEach(category => {
            const button = document.createElement('button');
            button.textContent = category.category;
            button.dataset.categoryId = category.id;
            button.classList.toggle('active', activeCategoryId === category.id);
            categoryFilterBar.appendChild(button);
        });

        categoryFilterBar.addEventListener('click', (e) => {
            if (e.target.tagName === 'BUTTON') {
                activeCategoryId = e.target.dataset.categoryId === 'all' ? 'all' : parseInt(e.target.dataset.categoryId, 10);
                // Re-render buttons to update active state
                renderCategories(); 
                // Filter and render products
                productSearch.value = ''; // Clear search
                const productsToDisplay = activeCategoryId === 'all' 
                    ? allProducts 
                    : allProducts.filter(p => p.category.id === activeCategoryId);
                renderProducts(productsToDisplay);
            }
        });
    };

    const renderCart = () => {
        cartTableBody.innerHTML = '';
        let totalItems = 0;
        let uniqueItems = 0;
        let totalAmount = 0;
        cart.forEach(item => {
            const row = cartItemRowTemplate.content.cloneNode(true);
            const tr = row.querySelector('tr');
            tr.dataset.productId = item.id;
            row.querySelector('.product-name').textContent = item.name;
            row.querySelector('.product-code').textContent = `(${item.code})`;
            row.querySelector('.product-price').textContent = formatCurrency(item.price);
            row.querySelector('.qty-value').textContent = item.quantity;
            row.querySelector('.product-subtotal').textContent = formatCurrency(item.price * item.quantity);
            cartTableBody.appendChild(row);
            totalItems += item.quantity;
            uniqueItems++;
            totalAmount += item.price * item.quantity;
        });
        // Update summaries
        summaryTotalItems.textContent = `${totalItems} (${uniqueItems})`;
        summaryTotal.textContent = formatCurrency(totalAmount);
        summaryTotalPayable.textContent = formatCurrency(totalAmount); // TODO: Add discount and tax logic
        checkoutButton.disabled = cart.length === 0;
    };

    // --- Logic Functions ---
  const addToCart = (productId) => {
        const product = allProducts.find(p => p.id === productId);
        if (!product) return;

        // Check if product has stock
        if (product.total_stock <= 0) {
            alert('Product is out of stock');
            return;
        }

        const existingItem = cart.find(item => item.id === productId);
        const totalQuantity = existingItem ? existingItem.quantity + 1 : 1;

        // Check if requested quantity exceeds stock
        if (totalQuantity > product.total_stock) {
            alert(`Insufficient stock. Available: ${product.total_stock}`);
            return;
        }

        if (existingItem) {
            existingItem.quantity++;
        } else {
            cart.push({ id: product.id, name: product.name, code: product.sku, price: product.price, quantity: 1 });
        }
        renderCart();
    };
    const updateCartQuantity = (productId, change) => {
        const itemIndex = cart.findIndex(item => item.id === productId);
        if (itemIndex === -1) return;

        const newQuantity = cart[itemIndex].quantity + change;
        const product = allProducts.find(p => p.id === productId);

        if (change > 0 && product && newQuantity > product.total_stock) {
            alert(`Insufficient stock. Available: ${product.total_stock}`);
            return;
        }

        cart[itemIndex].quantity = newQuantity;
        if (cart[itemIndex].quantity <= 0) {
            cart.splice(itemIndex, 1); // Remove item
        }
        renderCart();
    };
    const handleCheckout = async () => {
        if (cart.length === 0) return;
        const posWrapper = document.querySelector('.pos-wrapper');
        const tokoId = parseInt(posWrapper.dataset.tokoId, 10);
        const userId = parseInt(posWrapper.dataset.userId, 10);
        if (!tokoId || !userId) {
            alert('Error: Missing TokoID or UserID. Please login again.');
            return;
        }
        const payload = {
            toko_id: tokoId,
            user_id: userId,
            payment_method: 'Cash', // Hardcoded for now
            items: cart.map(item => ({ product_id: item.id, quantity: item.quantity }))
        };
        try {
            const response = await fetch('/api/pos/transactions', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Transaction failed');
            }
            alert('Transaksi berhasil!');
            cart = [];
            renderCart();
        } catch (error) {
            console.error('Checkout error:', error);
            alert(`Error: ${error.message}`);
        }
    };
    // --- Event Listeners ---
    productGrid.addEventListener('click', (e) => {
        const card = e.target.closest('.product-card');
        if (card) {
            addToCart(parseInt(card.dataset.productId, 10));
        }
    });
    cartTableBody.addEventListener('click', (e) => {
        const row = e.target.closest('tr');
        if (!row) return;
        const productId = parseInt(row.dataset.productId, 10);
        if (e.target.classList.contains('btn-qty-increase')) {
            updateCartQuantity(productId, 1);
        } else if (e.target.classList.contains('btn-qty-decrease')) {
            updateCartQuantity(productId, -1);
        } else if (e.target.classList.contains('btn-remove-item')) {
            updateCartQuantity(productId, -Infinity); // Remove item completely
        }
    });
    productSearch.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        
        // Start with the currently selected category's products
        let productsToFilter = activeCategoryId === 'all'
            ? allProducts
            : allProducts.filter(p => p.category.id === activeCategoryId);

        const filteredProducts = productsToFilter.filter(p =>
            p.name.toLowerCase().includes(searchTerm) ||
            p.sku.toLowerCase().includes(searchTerm)
        );
        renderProducts(filteredProducts);
    });
    checkoutButton.addEventListener('click', handleCheckout);
    cancelCartButton.addEventListener('click', () => {
        if (confirm('Are you sure you want to cancel the cart?')) {
            cart = [];
            renderCart();
        }
    });
    // --- Initialization ---
    const initializePos = async () => {
        try {
            const response = await fetch('/api/pos/products');
            if (!response.ok) throw new Error('Failed to fetch products');
            const products = await response.json();
            allProducts = products || [];

            // Extract unique categories
            const categoryMap = new Map();
            allProducts.forEach(p => {
                if (p.category && p.category.id) {
                    categoryMap.set(p.category.id, p.category);
                }
            });
            categories = Array.from(categoryMap.values()).sort((a, b) => a.category.localeCompare(b.category));

            renderCategories();
            renderProducts(allProducts);
            renderCart();
            updateClock();
            setInterval(updateClock, 1000);
        } catch (error) {
            console.error('Initialization failed:', error);
            alert('Could not initialize POS. Please check connection and refresh.');
        }
    };

    initializePos();
});