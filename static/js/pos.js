document.addEventListener('DOMContentLoaded', () => {
    // --- DOM Elements ---
    const productGrid = document.getElementById('product-grid');
    const cartTableBody = document.getElementById('cart-items-table-body');
    const productSearch = document.getElementById('product-search');
    const barcodeScanBtn = document.getElementById('barcode-scan-btn');
    const checkoutButton = document.querySelector('.btn-payment');
    const cancelCartButton = document.querySelector('.btn-cancel');
    const holdButton = document.querySelector('.btn-hold');
    const printOrderButton = document.querySelector('.btn-print-order');
    const clockElement = document.getElementById('header-time');
    const categoryFilterBar = document.getElementById('category-filter-bar');

    // Member elements
    const customerSelect = document.getElementById('customer-select');
    const addMemberBtn = document.getElementById('add-member-btn');
    const memberInfoBtn = document.getElementById('member-info-btn');
    const referenceNote = document.getElementById('reference-note');
    const memberSearchContainer = document.getElementById('member-search-container');
    const memberSearchInput = document.getElementById('member-search-input');
    const memberSearchResults = document.getElementById('member-search-results');

    // Member modal elements
    const addMemberModal = document.getElementById('add-member-modal');
    const memberModalClose = document.getElementById('member-modal-close');
    const newMemberName = document.getElementById('new-member-name');
    const newMemberPhone = document.getElementById('new-member-phone');
    const newMemberEmail = document.getElementById('new-member-email');
    const newMemberAddress = document.getElementById('new-member-address');
    const newMemberBirthday = document.getElementById('new-member-birthday');
    const newMemberGender = document.getElementById('new-member-gender');
    const newMemberTier = document.getElementById('new-member-tier');
    const newMemberNotes = document.getElementById('new-member-notes');
    const saveMemberBtn = document.getElementById('save-member-btn');
    const cancelMemberBtn = document.getElementById('cancel-member-btn');

    // Member info modal elements
    const memberInfoModal = document.getElementById('member-info-modal');
    const memberInfoModalClose = document.getElementById('member-info-modal-close');
    const closeMemberInfoBtn = document.getElementById('close-member-info-btn');

    // Member display elements
    const memberInfoDisplay = document.getElementById('member-info-display');
    const selectedMemberName = document.getElementById('selected-member-name');
    const selectedMemberPoints = document.getElementById('selected-member-points');
    const selectedMemberTier = document.getElementById('selected-member-tier');
    const summaryMemberDiscount = document.getElementById('summary-member-discount');
    const summaryPointsEarned = document.getElementById('summary-points-earned');
    const redeemPointsSection = document.getElementById('redeem-points-section');
    const redeemPointsInput = document.getElementById('redeem-points-input');
    const btnRedeemPoints = document.getElementById('btn-redeem-points');

    // Held transactions elements
    const heldTransactionsModal = document.getElementById('held-transactions-modal');
    const heldModalClose = document.getElementById('held-modal-close');
    const heldTransactionsList = document.getElementById('held-transactions-list');
    const noHeldTransactions = document.getElementById('no-held-transactions');
    const cancelHeldBtn = document.getElementById('cancel-held-btn');
    const viewHeldBtn = document.getElementById('view-held-btn');

    // Barcode scanner elements
    const barcodeScannerModal = document.getElementById('barcode-scanner-modal');
    const scannerModalClose = document.getElementById('scanner-modal-close');
    const stopScanBtn = document.getElementById('stop-scan-btn');
    const testScanBtn = document.getElementById('test-scan-btn');
    const scanResult = document.getElementById('scan-result');
    const qrReaderContainer = document.getElementById('qr-reader');

    // Summary elements
    const summaryTotalItems = document.getElementById('summary-total-items');
    const summaryTotal = document.getElementById('summary-total');
    const summaryDiscount = document.getElementById('summary-discount');
    const summaryTax = document.getElementById('summary-tax');
    const summaryTotalPayable = document.getElementById('summary-total-payable');

    // Payment modal elements
    const paymentModal = document.getElementById('payment-modal');
    const receiptModal = document.getElementById('receipt-modal');
    const paymentModalClose = document.getElementById('payment-modal-close');
    const receiptModalClose = document.getElementById('receipt-modal-close');
    const confirmPaymentBtn = document.getElementById('confirm-payment-btn');
    const cancelPaymentBtn = document.getElementById('cancel-payment-btn');
    const printReceiptBtn = document.getElementById('print-receipt-btn');
    const closeReceiptBtn = document.getElementById('close-receipt-btn');

    // Payment form elements
    const paymentMethodSelect = document.getElementById('payment-method-select');
    const cashAmountInput = document.getElementById('cash-amount');
    const customerNameInput = document.getElementById('customer-name');
    const cashPaymentSection = document.getElementById('cash-payment-section');

    // Templates
    const productCardTemplate = document.getElementById('product-card-template');
    const cartItemRowTemplate = document.getElementById('cart-item-row-template');

    // --- State ---
    let allProducts = [];
    let cart = [];
    let categories = [];
    let activeCategoryId = 'all';
    let barcodeBuffer = '';
    let currentTransaction = null;
    let html5QrCode = null;
    let isScanning = false;

    // Member state
    let members = [];
    let selectedMember = null;
    let isMemberSearchVisible = false;
    let memberSearchTimeout = null;

    // Held transactions state
    let heldTransactions = [];
    let heldTransactionCounter = 1;

    // --- Utility Functions ---
    const formatCurrency = (amount) =>
        new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(amount);

    // --- Member Functions ---
    const searchMembers = async (query) => {
        if (!query || query.length < 3) {
            memberSearchResults.innerHTML = '';
            return;
        }

        // Clear previous timeout
        if (memberSearchTimeout) {
            clearTimeout(memberSearchTimeout);
        }

        memberSearchTimeout = setTimeout(async () => {
            memberSearchResults.innerHTML = '<div class="member-search-loading">Searching...</div>';

            try {
                // Try API first with updated route
                const response = await fetch(`/api/member-search?q=${encodeURIComponent(query)}`);
                if (response.ok) {
                    const searchResults = await response.json();
                    displayMemberSearchResults(searchResults);
                } else {
                    // Fallback to local mock data
                    displayMemberSearchResults([]);
                }
            } catch (error) {
                console.log('Failed to search members:', error);
                memberSearchResults.innerHTML = '<div class="member-search-loading">No members found</div>';
            }
        }, 300);
    };

    const displayMemberSearchResults = (results) => {
        if (results.length === 0) {
            memberSearchResults.innerHTML = `
                <div class="member-search-result-item">
                    <div class="member-result-name">No members found</div>
                    <div class="member-result-details">Try different search terms</div>
                </div>
            `;
            return;
        }

        memberSearchResults.innerHTML = results.map(member => `
            <div class="member-search-result-item" onclick="selectMember(${member.id})">
                <div class="member-result-name">${member.name}</div>
                <div class="member-result-details">
                    <span>${member.phone}</span>
                    <span class="member-result-code">${member.memberCode}</span>
                    <span class="member-result-points">${member.totalPoints} pts</span>
                </div>
            </div>
        `).join('');
    };

    const selectMember = async (memberId) => {
        // Validate member ID
        if (!memberId || memberId === 'undefined' || memberId === 'null') {
            console.error('Invalid member ID:', memberId);
            alert('Invalid member ID');
            return;
        }

        try {
            const response = await fetch(`/api/members/${memberId}`);
            if (response.ok) {
                selectedMember = await response.json();
                displaySelectedMember();
                hideMemberSearch();
            } else {
                console.error('Failed to load member, status:', response.status);
                alert('Failed to load member details');
            }
        } catch (error) {
            console.error('Failed to load member:', error);
            alert('Failed to load member details');
        }
    };

    const displaySelectedMember = () => {
        if (!selectedMember) {
            memberInfoDisplay.style.display = 'none';
            memberInfoBtn.style.display = 'none';
            redeemPointsSection.style.display = 'none';
            return;
        }

        memberInfoDisplay.style.display = 'grid';
        memberInfoBtn.style.display = 'block';
        redeemPointsSection.style.display = 'block';

        selectedMemberName.textContent = selectedMember.name;
        selectedMemberPoints.textContent = `${selectedMember.totalPoints} pts`;
        selectedMemberTier.textContent = selectedMember.memberTier?.name || 'Bronze';

        updateMemberDiscount();
        calculatePointsEarned();
    };

    const updateMemberDiscount = () => {
        if (!selectedMember || !selectedMember.memberTier) {
            summaryMemberDiscount.textContent = 'Rp 0';
            return;
        }

        const totals = calculateTotals();
        const discountRate = selectedMember.memberTier.discountRate || 0;
        const memberDiscount = totals.subtotal * discountRate;

        summaryMemberDiscount.textContent = formatCurrency(memberDiscount);
        updateTotals();
    };

    const calculatePointsEarned = () => {
        if (!selectedMember) {
            summaryPointsEarned.textContent = '0 pts';
            return;
        }

        const totals = calculateTotals();
        const pointRate = selectedMember.memberTier?.pointRate || 1;
        const pointsEarned = Math.floor(totals.subtotal / 100) * pointRate; // 1 point per Rp 100

        summaryPointsEarned.textContent = `${pointsEarned} pts`;
    };

    const showMemberSearch = () => {
        memberSearchContainer.style.display = 'block';
        memberSearchInput.value = '';
        memberSearchInput.focus();
        isMemberSearchVisible = true;
    };

    const hideMemberSearch = () => {
        memberSearchContainer.style.display = 'none';
        memberSearchInput.value = '';
        memberSearchResults.innerHTML = '';
        isMemberSearchVisible = false;
    };

    const showAddMemberModal = async () => {
        // Load member tiers
        try {
            const response = await fetch('/api/member-tiers');
            if (response.ok) {
                const tiers = await response.json();
                newMemberTier.innerHTML = '<option value="">Default Tier (Auto-assigned)</option>';
                tiers.forEach(tier => {
                    newMemberTier.innerHTML += `<option value="${tier.id}">${tier.name}</option>`;
                });
            }
        } catch (error) {
            console.log('Failed to load member tiers:', error);
        }

        newMemberName.value = '';
        newMemberPhone.value = '';
        newMemberEmail.value = '';
        newMemberAddress.value = '';
        addMemberModal.style.display = 'block';
        newMemberName.focus();
    };

    const hideAddMemberModal = () => {
        addMemberModal.style.display = 'none';
    };

    const addMember = async () => {
        // Get all form values
        const name = newMemberName.value.trim();
        const phone = newMemberPhone.value.trim();
        const email = newMemberEmail.value.trim();
        const address = newMemberAddress.value.trim();
        const birthday = newMemberBirthday?.value.trim();
        const gender = newMemberGender?.value.trim();
        const tier = newMemberTier.value.trim();
        const notes = newMemberNotes?.value.trim();

        // Enhanced validation
        const errors = [];

        if (!name) {
            errors.push('Member name is required');
        } else if (name.length < 3) {
            errors.push('Member name must be at least 3 characters');
        }

        if (!phone) {
            errors.push('Phone number is required');
        } else if (!/^[\d\s\-\+\(\)]+$/.test(phone)) {
            errors.push('Please enter a valid phone number');
        } else if (phone.replace(/\D/g, '').length < 10) {
            errors.push('Phone number must be at least 10 digits');
        }

        if (email && ! /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
            errors.push('Please enter a valid email address');
        }

        if (errors.length > 0) {
            alert('Please fix the following errors:\n' + errors.join('\n'));
            return;
        }

        // Create member data object
        const newMemberData = {
            name: name,
            phone: phone,
            email: email || null,
            address: address || null
        };

        // Add optional fields if provided
        if (birthday) newMemberData.birthday = birthday;
        if (gender) newMemberData.gender = gender;
        if (tier) newMemberData.tierId = parseInt(tier);
        if (notes) newMemberData.notes = notes;

        try {
            // Show loading state
            const saveBtn = document.getElementById('save-member-btn');
            const originalText = saveBtn.innerHTML;
            saveBtn.innerHTML = '<i class="bi bi-hourglass-split"></i> Adding Member...';
            saveBtn.disabled = true;

            const response = await fetch('/api/members', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(newMemberData)
            });

            if (response.ok) {
                const savedMember = await response.json();
                console.log('Member saved to database:', savedMember);

                // Auto-select the new member if ID exists
                if (savedMember && savedMember.id) {
                    selectMember(savedMember.id);
                    hideAddMemberModal();
                    alert('Member added successfully!\nMember Code: ' + (savedMember.memberCode || 'Generated'));
                } else {
                    console.error('Invalid member response:', savedMember);
                    alert('Member added but failed to load member details');
                    hideAddMemberModal();
                }
            } else {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Unknown error');
            }
        } catch (error) {
            console.error('Failed to add member:', error);
            alert('Failed to add member: ' + error.message);
        } finally {
            // Restore button state
            const saveBtn = document.getElementById('save-member-btn');
            saveBtn.innerHTML = '<i class="bi bi-person-plus"></i> Add Member';
            saveBtn.disabled = false;
        }
    };

    const showMemberInfoModal = async () => {
        if (!selectedMember) return;

        // Populate member info
        document.getElementById('info-member-code').textContent = selectedMember.memberCode;
        document.getElementById('info-member-name').textContent = selectedMember.name;
        document.getElementById('info-member-phone').textContent = selectedMember.phone;
        document.getElementById('info-member-email').textContent = selectedMember.email || '-';
        document.getElementById('info-member-tier').textContent = selectedMember.memberTier?.name || 'Bronze';
        document.getElementById('info-member-points').textContent = selectedMember.totalPoints;
        document.getElementById('info-member-spent').textContent = formatCurrency(selectedMember.totalSpent);
        document.getElementById('info-member-joined').textContent = new Date(selectedMember.joinedDate).toLocaleDateString('id-ID');
        document.getElementById('info-member-last-visit').textContent = selectedMember.lastVisitDate ?
            new Date(selectedMember.lastVisitDate).toLocaleDateString('id-ID') : '-';

        // Load member transactions
        loadMemberTransactions();

        memberInfoModal.style.display = 'block';
    };

    const hideMemberInfoModal = () => {
        memberInfoModal.style.display = 'none';
    };

    const loadMemberTransactions = async () => {
        if (!selectedMember) return;

        try {
            const response = await fetch(`/api/member-transactions/${selectedMember.id}?limit=10`);
            if (response.ok) {
                const transactions = await response.json();
                displayMemberTransactions(transactions);
            }
        } catch (error) {
            console.log('Failed to load member transactions:', error);
        }
    };

    const displayMemberTransactions = (transactions) => {
        const transactionsList = document.getElementById('member-transactions-list');

        if (transactions.length === 0) {
            transactionsList.innerHTML = '<p>No transactions found</p>';
            return;
        }

        transactionsList.innerHTML = transactions.map(transaction => `
            <div class="transaction-item ${transaction.type.toLowerCase()}">
                <div class="transaction-details">
                    <div class="transaction-description">${transaction.description}</div>
                    <div class="transaction-date">${new Date(transaction.createdAt).toLocaleDateString('id-ID')}</div>
                </div>
                <div class="transaction-points ${transaction.points > 0 ? 'positive' : 'negative'}">
                    ${transaction.points > 0 ? '+' : ''}${transaction.points} pts
                </div>
            </div>
        `).join('');
    };

    const redeemPoints = async () => {
        if (!selectedMember) {
            alert('No member selected');
            return;
        }

        const pointsToRedeem = parseInt(redeemPointsInput.value) || 0;

        if (pointsToRedeem <= 0) {
            alert('Please enter valid points to redeem');
            return;
        }

        if (pointsToRedeem > selectedMember.totalPoints) {
            alert('Insufficient points');
            return;
        }

        try {
            const response = await fetch(`/api/member-redeem/${selectedMember.id}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    points: pointsToRedeem,
                    description: `Redeemed ${pointsToRedeem} points for discount`
                })
            });

            if (response.ok) {
                const redemption = await response.json();

                // Update member points
                selectedMember.totalPoints = redemption.balanceAfter;
                displaySelectedMember();

                // Add discount (simple conversion: 100 points = Rp 1000)
                const discountValue = pointsToRedeem * 10; // 1 point = Rp 10
                const currentDiscount = parseFloat(summaryDiscount.value) || 0;
                summaryDiscount.value = currentDiscount + discountValue;

                redeemPointsInput.value = 0;
                updateTotals();

                alert(`Successfully redeemed ${pointsToRedeem} points for ${formatCurrency(discountValue)} discount!`);
            } else {
                const errorData = await response.json();
                alert('Failed to redeem points: ' + (errorData.error || 'Unknown error'));
            }
        } catch (error) {
            console.log('Failed to redeem points:', error);
            alert('Failed to redeem points. Please try again.');
        }
    };

    // --- Held Transactions Functions ---
    const holdTransaction = async () => {
        if (cart.length === 0) {
            alert('Cart is empty. Nothing to hold.');
            return;
        }

        const totals = calculateTotals();
        const heldTransaction = {
            id: 'hold_' + heldTransactionCounter++,
            items: JSON.parse(JSON.stringify(cart)), // Deep copy to prevent reference issues
            totals: JSON.parse(JSON.stringify(totals)), // Deep copy totals object
            customer: selectedMember ? JSON.parse(JSON.stringify(selectedMember)) : null,
            referenceNote: referenceNote.value,
            date: new Date().toISOString(),
            humanDate: new Date().toLocaleString('id-ID')
        };

        console.log('Holding transaction:', heldTransaction);
        heldTransactions.push(heldTransaction);

        // Save to localStorage with error handling
        try {
            localStorage.setItem('pos_held_transactions', JSON.stringify(heldTransactions));
            console.log('Held transactions saved to localStorage:', heldTransactions.length);
        } catch (error) {
            console.error('Failed to save held transactions:', error);
            alert('Error saving held transaction. Please try again.');
            return;
        }

        // Clear current cart after successful save
        cart = [];
        selectedMember = null;
        displaySelectedMember();
        customerSelect.value = '';
        referenceNote.value = '';
        renderCart();

        alert('Transaction held successfully! You can resume it later.');
    };

    
    const loadHeldTransactions = () => {
        try {
            const saved = localStorage.getItem('pos_held_transactions');
            if (saved) {
                heldTransactions = JSON.parse(saved);
                console.log('Loaded held transactions from localStorage:', heldTransactions.length);
            } else {
                heldTransactions = [];
                console.log('No held transactions found in localStorage');
            }
        } catch (error) {
            console.error('Error loading held transactions:', error);
            heldTransactions = [];
            // Clear corrupted data
            localStorage.removeItem('pos_held_transactions');
        }
    };

    const showHeldTransactions = () => {
        renderHeldTransactionsList();
        heldTransactionsModal.style.display = 'block';
    };

    const hideHeldTransactions = () => {
        heldTransactionsModal.style.display = 'none';
    };

    const renderHeldTransactionsList = () => {
        if (heldTransactions.length === 0) {
            heldTransactionsList.style.display = 'none';
            noHeldTransactions.style.display = 'block';
            return;
        }

        heldTransactionsList.style.display = 'block';
        noHeldTransactions.style.display = 'none';

        heldTransactionsList.innerHTML = heldTransactions.map(hold => `
            <div class="held-transaction-item">
                <div class="held-transaction-info">
                    <div class="held-transaction-customer">${hold.customer ? hold.customer.name : 'Walk-in Client'}</div>
                    <div class="held-transaction-date">${hold.humanDate}</div>
                    <div class="held-transaction-total">${formatCurrency(hold.totals.totalPayable)} (${hold.totals.totalItems} items)</div>
                </div>
                <div class="held-transaction-actions">
                    <button class="btn-resume" onclick="resumeHeldTransaction('${hold.id}')">Resume</button>
                    <button class="btn-delete-held" onclick="deleteHeldTransaction('${hold.id}')">Delete</button>
                </div>
            </div>
        `).join('');
    };

    const resumeHeldTransaction = async (heldId) => {
        try {
            const heldTransaction = heldTransactions.find(h => h.id === heldId);
            if (!heldTransaction) {
                console.error('Held transaction not found:', heldId);
                return;
            }

            console.log('Resuming held transaction:', heldTransaction);

            // Clear current cart
            cart = [];

            // Deep copy items to prevent reference issues
            cart = JSON.parse(JSON.stringify(heldTransaction.items));

            // Restore member if exists
            if (heldTransaction.customer) {
                selectedMember = JSON.parse(JSON.stringify(heldTransaction.customer));
                displaySelectedMember();
            } else {
                selectedMember = null;
                displaySelectedMember();
            }
            customerSelect.value = '';

            // Restore reference note
            referenceNote.value = heldTransaction.referenceNote || '';

            // Update UI
            renderCart();

            // Remove from held transactions and save
            heldTransactions = heldTransactions.filter(h => h.id !== heldId);
            localStorage.setItem('pos_held_transactions', JSON.stringify(heldTransactions));
            console.log('Removed held transaction from storage');

            hideHeldTransactions();
            alert('Transaction resumed successfully!');
        } catch (error) {
            console.error('Error resuming transaction:', error);
            alert('Error resuming transaction. Please try again.');
        }
    };

    const deleteHeldTransaction = (heldId) => {
        if (!confirm('Are you sure you want to delete this held transaction?')) return;

        try {
            heldTransactions = heldTransactions.filter(h => h.id !== heldId);
            localStorage.setItem('pos_held_transactions', JSON.stringify(heldTransactions));
            console.log('Deleted held transaction:', heldId);
            console.log('Remaining held transactions:', heldTransactions.length);
            renderHeldTransactionsList();
        } catch (error) {
            console.error('Error deleting held transaction:', error);
            alert('Error deleting held transaction. Please try again.');
        }
    };

    // --- Print Order Function ---
    const printOrder = () => {
        if (cart.length === 0) {
            alert('Cart is empty. Nothing to print.');
            return;
        }

        const totals = calculateTotals();
        const printWindow = window.open('', '_blank');

        let printContent = `
            <!DOCTYPE html>
            <html>
            <head>
                <title>Order Receipt</title>
                <style>
                    body {
                        font-family: 'Courier New', monospace;
                        margin: 20px;
                        font-size: 12px;
                    }
                    .header {
                        text-align: center;
                        margin-bottom: 20px;
                    }
                    .info {
                        margin-bottom: 15px;
                    }
                    .items {
                        margin-bottom: 15px;
                    }
                    .item {
                        margin: 5px 0;
                    }
                    .totals {
                        border-top: 1px dashed #000;
                        padding-top: 10px;
                        margin-bottom: 15px;
                    }
                    .total {
                        font-weight: bold;
                    }
                    .footer {
                        text-align: center;
                        margin-top: 20px;
                    }
                </style>
            </head>
            <body>
                <div class="header">
                    <h2>ORDER RECEIPT</h2>
                    <p>================================</p>
                </div>

                <div class="info">
                    <p>Date: ${new Date().toLocaleString('id-ID')}</p>
                    <p>Customer: ${selectedMember ? selectedMember.name : 'Walk-in Client'}</p>
                    <p>Ref: ${referenceNote.value || '-'}</p>
                    <p>================================</p>
                </div>

                <div class="items">
                    <p><strong>Items:</strong></p>
        `;

        cart.forEach(item => {
            printContent += `
                <div class="item">
                    ${item.name} (${item.quantity}x) = ${formatCurrency(item.price * item.quantity)}
                </div>
            `;
        });

        printContent += `
                </div>

                <div class="totals">
                    <p>Total Items: ${totals.totalItems}</p>
                    <p>Subtotal: ${formatCurrency(totals.subtotal)}</p>
                    <p>Discount: ${formatCurrency(totals.discount)}</p>
                    <p>Tax: ${formatCurrency(totals.tax)}</p>
                    <p class="total">Total: ${formatCurrency(totals.totalPayable)}</p>
                </div>

                <div class="footer">
                    <p>================================</p>
                    <p>Thank you for your order!</p>
                    <p>Please proceed to payment counter.</p>
                </div>
            </body>
            </html>
        `;

        printWindow.document.write(printContent);
        printWindow.document.close();
        printWindow.print();
    };

    const updateClock = () => {
        if (clockElement) {
            clockElement.textContent = new Date().toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'medium' });
        }
    };

    const calculateTotals = () => {
        let totalItems = 0;
        let subtotal = 0;

        cart.forEach(item => {
            totalItems += item.quantity;
            subtotal += item.price * item.quantity;
        });

        const discount = parseFloat(summaryDiscount.value) || 0;
        const tax = subtotal * 0.1; // 10% tax
        const totalPayable = Math.max(0, subtotal - discount + tax);

        return {
            totalItems,
            uniqueItems: cart.length,
            subtotal,
            discount,
            tax,
            totalPayable
        };
    };

    // --- Barcode Scanning ---
    const handleBarcodeInput = (e) => {
        if (e.key === 'Enter' && barcodeBuffer.length > 0) {
            // Search for product by barcode/SKU
            const product = allProducts.find(p => p.sku === barcodeBuffer);

            if (product) {
                addToCart(product.id);
                productSearch.value = '';
            } else {
                alert(`Product with barcode ${barcodeBuffer} not found`);
            }

            barcodeBuffer = '';
            e.preventDefault();
        } else if (e.key.length === 1) {
            barcodeBuffer += e.key;
        }
    };

    const startBarcodeScanner = async () => {
        try {
            console.log('Starting barcode scanner...');
            scanResult.textContent = 'Initializing scanner...';

            // Clear previous scanner
            qrReaderContainer.innerHTML = '';

            // Create new Html5Qrcode instance
            html5QrCode = new Html5Qrcode("qr-reader");
            isScanning = true;

            // Start scanning
            await html5QrCode.start(
                { facingMode: "environment" },
                {
                    fps: 10,
                    qrbox: { width: 250, height: 250 }
                },
                (decodedText, decodedResult) => {
                    console.log('Raw barcode detected:', decodedText);
                    console.log('Decoded result:', decodedResult);
                    handleScanResult(decodedText, decodedResult);
                },
                (errorMessage) => {
                    // Ignore scan errors for performance
                    // console.log('Scan error (ignored):', errorMessage);
                }
            );

            scanResult.textContent = 'Scanner ready! Point camera at barcode to scan products.';
            console.log('Scanner started successfully');

        } catch (error) {
            console.error('Error starting scanner:', error);
            scanResult.textContent = `Error: ${error.message}`;
            isScanning = false;

            // Provide user-friendly error messages
            if (error.name === 'NotAllowedError') {
                alert('Camera permission denied. Please allow camera access and try again.');
            } else if (error.name === 'NotFoundError') {
                alert('No camera found. Please ensure you have a camera connected.');
            } else {
                alert(`Could not start scanner: ${error.message}`);
            }
        }
    };

    const stopBarcodeScanner = () => {
        console.log('Stopping scanner...');

        try {
            if (html5QrCode && isScanning) {
                html5QrCode.stop().then(() => {
                    console.log('Scanner stopped successfully');
                    html5QrCode.clear();
                    qrReaderContainer.innerHTML = '';
                    isScanning = false;
                    html5QrCode = null;
                }).catch((error) => {
                    console.warn('Error stopping scanner:', error);
                    qrReaderContainer.innerHTML = '';
                    isScanning = false;
                    html5QrCode = null;
                });
            } else {
                qrReaderContainer.innerHTML = '';
                isScanning = false;
                html5QrCode = null;
            }

            scanResult.textContent = '';
            barcodeScannerModal.style.display = 'none';

        } catch (error) {
            console.error('Error stopping scanner:', error);
            // Force close even if there's an error
            qrReaderContainer.innerHTML = '';
            isScanning = false;
            html5QrCode = null;
            scanResult.textContent = '';
            barcodeScannerModal.style.display = 'none';
        }
    };

    const testBarcodeDetection = () => {
        console.log('=== TESTING BARCODE DETECTION ===');
        console.log('Total products available:', allProducts.length);

        if (allProducts.length === 0) {
            console.log('ERROR: No products loaded!');
            scanResult.textContent = 'ERROR: No products loaded!';
            return;
        }

        // Show first 5 products for debugging
        console.log('First 5 products:');
        allProducts.slice(0, 5).forEach((product, index) => {
            console.log(`${index + 1}. ID: ${product.id}, Name: "${product.name}", SKU: "${product.sku}"`);
        });

        // Test with common barcode patterns
        const testCodes = [
            '123456789',
            '987654321',
            'SKU001',
            'TEST001',
            allProducts[0]?.sku,
            allProducts[0]?.name
        ].filter(Boolean);

        console.log('Testing with codes:', testCodes);

        testCodes.forEach(code => {
            const product = allProducts.find(p =>
                p.sku === code || p.name === code ||
                p.name.toLowerCase().includes(code.toLowerCase()) ||
                p.sku.includes(code)
            );
            console.log(`Testing code "${code}":`, product ? `FOUND - ${product.name}` : 'NOT FOUND');
        });

        scanResult.innerHTML = `
            <strong>Test Results:</strong><br>
            • Total products: ${allProducts.length}<br>
            • Test codes checked: ${testCodes.length}<br>
            • Check console for details<br>
            <small>Console shows available products and search results</small>
        `;
    };

    
    const handleScanResult = (decodedText, decodedResult) => {
        console.log('Barcode detected:', decodedText);
        console.log('Full result:', decodedResult);

        // Show immediate feedback
        scanResult.textContent = `Scanned: ${decodedText}`;

        // Validate barcode
        if (!decodedText || decodedText.trim().length < 1) {
            scanResult.textContent = 'Invalid barcode detected';
            return;
        }

        const barcodeText = decodedText.trim();
        console.log('Clean barcode text:', barcodeText);
        console.log('Available products:', allProducts.map(p => ({
            id: p.id,
            name: p.name,
            sku: p.sku
        })));

        // Multiple search strategies
        let product = null;

        // Strategy 1: Exact SKU match
        product = allProducts.find(p => p.sku === barcodeText);
        if (product) {
            console.log('Found by exact SKU:', product);
        }

        // Strategy 2: Exact name match
        if (!product) {
            product = allProducts.find(p => p.name === barcodeText);
            if (product) {
                console.log('Found by exact name:', product);
            }
        }

        // Strategy 3: Partial SKU match (if barcode contains SKU)
        if (!product) {
            product = allProducts.find(p =>
                barcodeText.includes(p.sku) || p.sku.includes(barcodeText)
            );
            if (product) {
                console.log('Found by partial SKU match:', product);
            }
        }

        // Strategy 4: Partial name match (case insensitive)
        if (!product) {
            product = allProducts.find(p =>
                p.name.toLowerCase().includes(barcodeText.toLowerCase()) ||
                barcodeText.toLowerCase().includes(p.name.toLowerCase())
            );
            if (product) {
                console.log('Found by partial name match:', product);
            }
        }

        // Strategy 5: Try to extract numeric part and search
        if (!product && /\d/.test(barcodeText)) {
            const numericPart = barcodeText.replace(/\D/g, '');
            console.log('Trying numeric search with:', numericPart);

            if (numericPart.length > 0) {
                product = allProducts.find(p =>
                    p.sku.includes(numericPart) ||
                    p.name.includes(numericPart) ||
                    p.id.toString() === numericPart
                );
                if (product) {
                    console.log('Found by numeric match:', product);
                }
            }
        }

        if (product) {
            console.log('SUCCESS - Product found:', product);
            addToCart(product.id);
            scanResult.innerHTML = `
                ✓ Found: ${product.name}<br>
                <small>SKU: ${product.sku} | Added to cart!</small><br>
                <small>Scanner ready for next scan...</small>
            `;

            // Remove auto-close - scanner stays open for continuous scanning
            // Reset message after delay to show scanner is ready
            setTimeout(() => {
                scanResult.textContent = 'Scanner ready! Point camera at next barcode.';
            }, 3000);

        } else {
            console.log('FAILED - Product not found for barcode:', barcodeText);
            scanResult.innerHTML = `
                ✗ Product not found<br>
                <small>Barcode: "${barcodeText}"</small><br>
                <small>Try a different product or check barcode</small>
            `;

            // Reset scanner message after delay
            setTimeout(() => {
                scanResult.textContent = 'Scanner ready! Point camera at next barcode.';
            }, 3000);
        }
    };

    // --- Render Functions ---
    const renderProducts = (productsToRender) => {
        productGrid.innerHTML = '';
        productsToRender.forEach(product => {
            const card = productCardTemplate.content.cloneNode(true);
            const cardElement = card.querySelector('.product-card');
            cardElement.dataset.productId = product.id;
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
                renderCategories();
                productSearch.value = '';
                const productsToDisplay = activeCategoryId === 'all'
                    ? allProducts
                    : allProducts.filter(p => p.category && p.category.id === activeCategoryId);
                renderProducts(productsToDisplay);
            }
        });
    };

    const renderCart = () => {
        cartTableBody.innerHTML = '';
        const totals = calculateTotals();

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
        });

        // Update summaries
        summaryTotalItems.textContent = `${totals.totalItems} (${totals.uniqueItems})`;
        summaryTotal.textContent = formatCurrency(totals.subtotal);
        summaryTax.textContent = formatCurrency(totals.tax);
        summaryTotalPayable.textContent = formatCurrency(totals.totalPayable);
        checkoutButton.disabled = cart.length === 0;
    };

    // --- Logic Functions ---
    const addToCart = (productId) => {
        const product = allProducts.find(p => p.id === productId);
        if (!product) return;

        if (product.total_stock <= 0) {
            alert('Product is out of stock');
            return;
        }

        const existingItem = cart.find(item => item.id === productId);
        const totalQuantity = existingItem ? existingItem.quantity + 1 : 1;

        if (totalQuantity > product.total_stock) {
            alert(`Insufficient stock. Available: ${product.total_stock}`);
            return;
        }

        if (existingItem) {
            existingItem.quantity++;
        } else {
            cart.push({
                id: product.id,
                name: product.name,
                code: product.sku,
                price: product.price,
                quantity: 1
            });
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
            cart.splice(itemIndex, 1);
        }
        renderCart();
    };

    const showPaymentModal = () => {
        if (cart.length === 0) return;

        const totals = calculateTotals();

        // Update payment modal with current totals
        document.getElementById('payment-total-items').textContent = totals.totalItems;
        document.getElementById('payment-subtotal').textContent = formatCurrency(totals.subtotal);
        document.getElementById('payment-discount').textContent = formatCurrency(totals.discount);
        document.getElementById('payment-tax').textContent = formatCurrency(totals.tax);
        document.getElementById('payment-total-payable').textContent = formatCurrency(totals.totalPayable);

        // Reset form
        paymentMethodSelect.value = 'Cash';
        cashAmountInput.value = '';
        customerNameInput.value = 'Walk-in Customer';
        updateChangeAmount();

        paymentModal.style.display = 'block';
    };

    const updateChangeAmount = () => {
        const totals = calculateTotals();
        const cashAmount = parseFloat(cashAmountInput.value) || 0;
        const change = Math.max(0, cashAmount - totals.totalPayable);
        document.getElementById('change-amount').textContent = formatCurrency(change);
    };

    const processPayment = async () => {
        const totals = calculateTotals();
        const paymentMethod = paymentMethodSelect.value;
        const customerName = customerNameInput.value || 'Walk-in Customer';

        if (paymentMethod === 'Cash') {
            const cashAmount = parseFloat(cashAmountInput.value) || 0;
            if (cashAmount < totals.totalPayable) {
                alert('Insufficient cash amount');
                return;
            }
        }

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
            payment_method: paymentMethod,
            customer_name: customerName,
            cash_amount: paymentMethod === 'Cash' ? parseFloat(cashAmountInput.value) || 0 : 0,
            discount: totals.discount,
            tax: totals.tax,
            total_payable: totals.totalPayable,
            items: cart.map(item => ({
                product_id: item.id,
                quantity: item.quantity,
                price: item.price
            }))
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

            const transactionData = await response.json();
            currentTransaction = { ...transactionData, cart: [...cart] };

            paymentModal.style.display = 'none';
            showReceiptModal();

            cart = [];
            renderCart();
        } catch (error) {
            console.error('Checkout error:', error);
            alert(`Error: ${error.message}`);
        }
    };

    const showReceiptModal = () => {
        if (!currentTransaction) return;

        const receiptContent = document.getElementById('receipt-content');
        const totals = calculateTotals();

        let receiptHTML = `
            <div class="receipt-header">
                <h2>TOKO KELONTONG</h2>
                <p>Jl. Contoh No. 123</p>
                <p>Telp: (021) 1234567</p>
                <p>================================</p>
            </div>
            <div class="receipt-info">
                <p>No: ${currentTransaction.id || 'TRX-' + Date.now()}</p>
                <p>Date: ${new Date().toLocaleString('id-ID')}</p>
                <p>Kasir: Admin</p>
                <p>Customer: ${currentTransaction.customer_name || 'Walk-in Customer'}</p>
                <p>================================</p>
            </div>
            <div class="receipt-items">
        `;

        currentTransaction.cart.forEach(item => {
            receiptHTML += `
                <div class="receipt-item">
                    <span>${item.name} (${item.quantity}x)</span>
                    <span>${formatCurrency(item.price * item.quantity)}</span>
                </div>
            `;
        });

        receiptHTML += `
            </div>
            <div class="receipt-summary">
                <div class="receipt-item">
                    <span>Subtotal:</span>
                    <span>${formatCurrency(totals.subtotal)}</span>
                </div>
                <div class="receipt-item">
                    <span>Discount:</span>
                    <span>${formatCurrency(totals.discount)}</span>
                </div>
                <div class="receipt-item">
                    <span>Tax (10%):</span>
                    <span>${formatCurrency(totals.tax)}</span>
                </div>
                <div class="receipt-item" style="font-weight: bold; font-size: 1.1em;">
                    <span>TOTAL:</span>
                    <span>${formatCurrency(totals.totalPayable)}</span>
                </div>
        `;

        if (currentTransaction.payment_method === 'Cash' && currentTransaction.cash_amount > 0) {
            const change = currentTransaction.cash_amount - totals.totalPayable;
            receiptHTML += `
                <div class="receipt-item">
                    <span>Cash:</span>
                    <span>${formatCurrency(currentTransaction.cash_amount)}</span>
                </div>
                <div class="receipt-item">
                    <span>Change:</span>
                    <span>${formatCurrency(change)}</span>
                </div>
            `;
        }

        receiptHTML += `
            </div>
            <div class="receipt-footer">
                <p>================================</p>
                <p>Terima Kasih</p>
                <p>Selamat Berbelanja Kembali</p>
            </div>
        `;

        receiptContent.innerHTML = receiptHTML;
        receiptModal.style.display = 'block';
    };

    const printReceipt = () => {
        const receiptContent = document.getElementById('receipt-content').innerHTML;
        const printWindow = window.open('', '_blank');
        printWindow.document.write(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>Receipt</title>
                <style>
                    body { font-family: 'Courier New', monospace; margin: 20px; }
                    .receipt-header { text-align: center; margin-bottom: 20px; }
                    .receipt-info { margin-bottom: 15px; }
                    .receipt-items { margin: 15px 0; }
                    .receipt-item { display: flex; justify-content: space-between; margin: 5px 0; }
                    .receipt-summary { border-top: 1px dashed #333; padding-top: 10px; margin-top: 15px; }
                    .receipt-footer { text-align: center; margin-top: 20px; border-top: 1px dashed #333; padding-top: 10px; }
                </style>
            </head>
            <body>${receiptContent}</body>
            </html>
        `);
        printWindow.document.close();
        printWindow.print();
        printWindow.close();
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
            updateCartQuantity(productId, -Infinity);
        }
    });

    productSearch.addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();

        let productsToFilter = activeCategoryId === 'all'
            ? allProducts
            : allProducts.filter(p => p.category && p.category.id === activeCategoryId);

        const filteredProducts = productsToFilter.filter(p =>
            p.name.toLowerCase().includes(searchTerm) ||
            (p.sku && p.sku.toLowerCase().includes(searchTerm))
        );
        renderProducts(filteredProducts);
    });

    // Barcode scanning
    productSearch.addEventListener('keypress', handleBarcodeInput);
    barcodeScanBtn.addEventListener('click', () => {
        barcodeScannerModal.style.display = 'block';
        startBarcodeScanner();
    });

    // Payment modal events
    checkoutButton.addEventListener('click', showPaymentModal);
    confirmPaymentBtn.addEventListener('click', processPayment);
    cancelPaymentBtn.addEventListener('click', () => {
        paymentModal.style.display = 'none';
    });

    paymentMethodSelect.addEventListener('change', (e) => {
        cashPaymentSection.style.display = e.target.value === 'Cash' ? 'block' : 'none';
    });

    cashAmountInput.addEventListener('input', updateChangeAmount);

    // Receipt modal events
    printReceiptBtn.addEventListener('click', printReceipt);
    closeReceiptBtn.addEventListener('click', () => {
        receiptModal.style.display = 'none';
        currentTransaction = null;
    });

    // Modal close events
    paymentModalClose.addEventListener('click', () => {
        paymentModal.style.display = 'none';
    });

    receiptModalClose.addEventListener('click', () => {
        receiptModal.style.display = 'none';
        currentTransaction = null;
    });

    // Scanner modal events - simple and direct
    scannerModalClose.addEventListener('click', () => {
        console.log('Close button clicked');
        stopBarcodeScanner();
    });

    stopScanBtn.addEventListener('click', () => {
        console.log('Stop button clicked');
        stopBarcodeScanner();
    });

    testScanBtn.addEventListener('click', () => {
        console.log('Test button clicked');
        testBarcodeDetection();
    });

    window.addEventListener('click', (e) => {
        if (e.target === paymentModal) {
            paymentModal.style.display = 'none';
        }
        if (e.target === receiptModal) {
            receiptModal.style.display = 'none';
            currentTransaction = null;
        }
        if (e.target === barcodeScannerModal) {
            stopBarcodeScanner();
        }
        if (e.target === addCustomerModal) {
            hideAddCustomerModal();
        }
        if (e.target === heldTransactionsModal) {
            hideHeldTransactions();
        }
    });

    // ESC key listener for closing modals
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            if (barcodeScannerModal.style.display === 'block') {
                stopBarcodeScanner();
            } else if (addCustomerModal.style.display === 'block') {
                hideAddCustomerModal();
            } else if (heldTransactionsModal.style.display === 'block') {
                hideHeldTransactions();
            } else if (paymentModal.style.display === 'block') {
                paymentModal.style.display = 'none';
            } else if (receiptModal.style.display === 'block') {
                receiptModal.style.display = 'none';
                currentTransaction = null;
            }
        }
    });

    cancelCartButton.addEventListener('click', () => {
        if (confirm('Are you sure you want to cancel the cart?')) {
            cart = [];
            renderCart();
        }
    });

    summaryDiscount.addEventListener('input', renderCart);

    // Member management event listeners
    customerSelect.addEventListener('change', (e) => {
        const selection = e.target.value;
        if (selection === '') {
            selectedMember = null;
            displaySelectedMember();
        } else if (selection === 'search') {
            showMemberSearch();
        }
    });

    addMemberBtn.addEventListener('click', showAddMemberModal);
    memberModalClose.addEventListener('click', hideAddMemberModal);
    saveMemberBtn.addEventListener('click', addMember);
    cancelMemberBtn.addEventListener('click', hideAddMemberModal);

    // Member info modal event listeners
    memberInfoBtn.addEventListener('click', showMemberInfoModal);
    memberInfoModalClose.addEventListener('click', hideMemberInfoModal);
    closeMemberInfoBtn.addEventListener('click', hideMemberInfoModal);

    // Member search event listeners
    memberSearchInput.addEventListener('input', (e) => {
        searchMembers(e.target.value);
    });

    memberSearchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Escape') {
            hideMemberSearch();
            customerSelect.value = selectedMember ? selectedMember.id : '';
        }
    });

    // Redeem points event listeners
    btnRedeemPoints.addEventListener('click', redeemPoints);

    // Make selectMember globally available for inline onclick handlers
    window.selectMember = selectMember;

    // Held transactions event listeners
    holdButton.addEventListener('click', holdTransaction);
    viewHeldBtn.addEventListener('click', showHeldTransactions);
    heldModalClose.addEventListener('click', hideHeldTransactions);
    cancelHeldBtn.addEventListener('click', hideHeldTransactions);

    // Print order event listener
    printOrderButton.addEventListener('click', printOrder);

    // Global functions for inline onclick handlers
    window.resumeHeldTransaction = resumeHeldTransaction;
    window.deleteHeldTransaction = deleteHeldTransaction;

    // --- Initialization ---
    const initializePos = async () => {
        try {
            const response = await fetch('/api/pos/products');
            if (!response.ok) throw new Error('Failed to fetch products');
            const products = await response.json();
            allProducts = products || [];

            const categoryMap = new Map();
            allProducts.forEach(p => {
                if (p.category && p.category.id) {
                    categoryMap.set(p.category.id, p.category);
                }
            });
            categories = Array.from(categoryMap.values()).sort((a, b) => a.category.localeCompare(b.category));

            // Initialize held transactions
            loadHeldTransactions();

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