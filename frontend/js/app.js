const PAGE_SIZE = 6;

const state = {
  token: localStorage.getItem("token") || "",
  username: localStorage.getItem("username") || "",
  currentProductId: null,
  catalogPage: 1,
};

const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

function showToast(msg, type = "success") {
  const el = $("#toast");
  el.textContent = msg;
  el.className = `toast ${type}`;
  setTimeout(() => el.classList.add("hidden"), 3000);
}

async function api(path, options = {}) {
  const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
  if (state.token) headers.Authorization = `Bearer ${state.token}`;

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  const text = await res.text();
  let data = null;
  if (text) {
    try {
      data = JSON.parse(text);
    } catch {
      data = text;
    }
  }

  if (!res.ok) {
    const err = new Error(data?.error || res.statusText);
    err.status = res.status;
    err.data = data;
    throw err;
  }
  return { status: res.status, data };
}

function showView(name) {
  $$(".view").forEach((v) => v.classList.add("hidden"));
  const el = document.getElementById(`view-${name}`);
  if (el) el.classList.remove("hidden");
}

function updateAuthUI() {
  const loggedIn = !!state.token;
  $("#nav-auth").classList.toggle("hidden", loggedIn);
  $("#nav-user").classList.toggle("hidden", !loggedIn);
  $("#bottom-nav").classList.toggle("hidden", !loggedIn);
  $("#username-label").textContent = state.username;
  $("#review-form").classList.toggle("hidden", !loggedIn);

  if (loggedIn) {
    showView("catalog");
    loadCatalog(1, true);
  } else {
    showView("login");
  }
}

function setAuth(token, username) {
  state.token = token;
  state.username = username;
  localStorage.setItem("token", token);
  localStorage.setItem("username", username);
  updateAuthUI();
}

function logout() {
  state.token = "";
  state.username = "";
  localStorage.removeItem("token");
  localStorage.removeItem("username");
  updateAuthUI();
}

$("#login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  $("#login-error").textContent = "";
  const fd = new FormData(e.target);
  try {
    const { data } = await api("/auth/login", {
      method: "POST",
      body: JSON.stringify({
        username: fd.get("username"),
        password: fd.get("password"),
      }),
    });
    setAuth(data.token, fd.get("username"));
    showToast("Welcome!");
  } catch (err) {
    $("#login-error").textContent = err.data?.error || err.message;
  }
});

$("#register-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  $("#register-error").textContent = "";
  const fd = new FormData(e.target);
  try {
    await api("/auth/register", {
      method: "POST",
      body: JSON.stringify({
        username: fd.get("username"),
        password: fd.get("password"),
      }),
    });
    const { data } = await api("/auth/login", {
      method: "POST",
      body: JSON.stringify({
        username: fd.get("username"),
        password: fd.get("password"),
      }),
    });
    setAuth(data.token, fd.get("username"));
    showToast("Account created");
  } catch (err) {
    $("#register-error").textContent = err.data?.error || err.message;
  }
});

$("#logout-btn").addEventListener("click", logout);

document.body.addEventListener("click", (e) => {
  const view = e.target.closest("[data-view]")?.dataset?.view;
  if (view) {
    showView(view);
    if (view === "catalog") loadCatalog(1, true);
    if (view === "orders") loadMyOrders();
    if (view === "admin") loadAdmin();
  }
});

$("#refresh-products").addEventListener("click", () => loadCatalog(1, true));
$("#filter-brand").addEventListener("change", () => loadCatalog(1));
$("#filter-category").addEventListener("change", () => loadCatalog(1));

$("#pagination-prev").addEventListener("click", () => {
  if (state.catalogPage > 1) loadCatalog(state.catalogPage - 1);
});

$("#pagination-next").addEventListener("click", () => {
  loadCatalog(state.catalogPage + 1);
});

async function loadFilters() {
  const [{ data: brands }, { data: categories }] = await Promise.all([
    api("/brands"),
    api("/categories"),
  ]);

  const brandSel = $("#filter-brand");
  const catSel = $("#filter-category");
  brandSel.innerHTML = '<option value="">All brands</option>';
  catSel.innerHTML = '<option value="">All categories</option>';

  brands.forEach((b) => {
    brandSel.innerHTML += `<option value="${b.id}">${b.name}</option>`;
  });
  categories.forEach((c) => {
    catSel.innerHTML += `<option value="${c.id}">${c.name}</option>`;
  });

  return { brands, categories };
}

function updatePagination(itemsOnPage) {
  const nav = $("#catalog-pagination");
  const hasPrev = state.catalogPage > 1;
  const hasNext = itemsOnPage === PAGE_SIZE;

  nav.classList.toggle("hidden", !hasPrev && !hasNext);
  $("#pagination-prev").disabled = !hasPrev;
  $("#pagination-next").disabled = !hasNext;
  $("#pagination-info").textContent = `Page ${state.catalogPage}`;
}

async function loadCatalog(page = state.catalogPage, reloadFilters = false) {
  try {
    if (reloadFilters) await loadFilters();

    state.catalogPage = page;
    const brand = $("#filter-brand").value;
    const category = $("#filter-category").value;

    let qs = `?page=${page}&limit=${PAGE_SIZE}`;
    if (brand) qs += `&brand_id=${brand}`;
    if (category) qs += `&category_id=${category}`;

    const { data: products } = await api(`/products${qs}`);
    const grid = $("#products-grid");
    grid.innerHTML = "";

    if (!products.length) {
      if (page > 1) {
        loadCatalog(page - 1);
        return;
      }
      grid.innerHTML = '<p class="meta">No products yet. Add some in Admin.</p>';
      updatePagination(0);
      return;
    }

    products.forEach((p) => {
      const card = document.createElement("article");
      card.className = "product-card";
      card.innerHTML = `
        <h3>${escapeHtml(p.name)}</h3>
        <p class="price">$${p.price.toFixed(2)}</p>
        <p class="meta">${escapeHtml(p.brand?.name || "")} · ${escapeHtml(p.category?.name || "")} · stock: ${p.stock}</p>
      `;
      card.addEventListener("click", () => openProduct(p.id));
      grid.appendChild(card);
    });

    updatePagination(products.length);
  } catch (err) {
    showToast(err.data?.error || err.message, "error");
  }
}

function escapeHtml(s) {
  const d = document.createElement("div");
  d.textContent = s ?? "";
  return d.innerHTML;
}

async function openProduct(id) {
  state.currentProductId = id;
  showView("product");

  const [{ data: product }, { data: reviews }] = await Promise.all([
    api(`/products/${id}`),
    api(`/products/${id}/reviews`),
  ]);

  $("#product-detail").innerHTML = `
    <div class="product-detail">
      <h2>${escapeHtml(product.name)}</h2>
      <p class="price">$${product.price.toFixed(2)}</p>
      <p class="meta">${escapeHtml(product.brand?.name || "")} · ${escapeHtml(product.category?.name || "")}</p>
      <p>In stock: ${product.stock} units</p>
    </div>
  `;

  const qtyInput = $('#order-form input[name="quantity"]');
  qtyInput.max = Math.max(product.stock, 1);
  qtyInput.value = product.stock > 0 ? 1 : 0;
  $("#order-card").classList.toggle("hidden", product.stock <= 0);

  const list = $("#reviews-list");
  list.innerHTML = reviews.length
    ? reviews
        .map(
          (r) => `
      <div class="review-item">
        <strong>${escapeHtml(r.user?.username || "user")}</strong>
        <span class="rating">${"★".repeat(r.rating)}${"☆".repeat(5 - r.rating)}</span>
        <span class="sentiment ${r.sentiment}">${r.sentiment || ""}</span>
        <p>${escapeHtml(r.comment || "")}</p>
      </div>
    `
        )
        .join("")
    : '<p class="meta">No reviews yet</p>';
}

$("#order-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const qty = Number(new FormData(e.target).get("quantity"));
  if (!state.currentProductId || qty < 1) return;

  try {
    await api("/orders", {
      method: "POST",
      body: JSON.stringify({
        items: [{ product_id: state.currentProductId, quantity: qty }],
      }),
    });
    showToast("Order placed");
    showView("orders");
    loadMyOrders();
  } catch (err) {
    showToast(err.data?.error || err.message, "error");
  }
});

async function loadMyOrders() {
  try {
    const { data: orders } = await api("/orders");
    const container = $("#orders-list");

    if (!orders.length) {
      container.innerHTML = '<p class="meta">You have no orders yet.</p>';
      return;
    }

    container.innerHTML = orders
      .map((order) => {
        const date = new Date(order.created_at).toLocaleString();
        const items = (order.items || [])
          .map(
            (item) => `
          <li>
            <span>${escapeHtml(item.product?.name || "Product")} × ${item.quantity}</span>
            <span>$${(item.price * item.quantity).toFixed(2)}</span>
          </li>`
          )
          .join("");

        return `
        <article class="order-card">
          <div class="order-card__header">
            <strong>Order #${order.id}</strong>
            <span class="order-card__status ${escapeHtml(order.status)}">${escapeHtml(order.status)}</span>
          </div>
          <p class="meta">${date}</p>
          <ul class="order-card__items">${items}</ul>
          <p class="order-card__total">Total: $${Number(order.total).toFixed(2)}</p>
        </article>`;
      })
      .join("");
  } catch (err) {
    showToast(err.data?.error || err.message, "error");
  }
}

$("#review-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  try {
    await api(`/products/${state.currentProductId}/reviews`, {
      method: "POST",
      body: JSON.stringify({
        rating: Number(fd.get("rating")),
        comment: fd.get("comment"),
      }),
    });
    e.target.reset();
    showToast("Review submitted");
    openProduct(state.currentProductId);
  } catch (err) {
    showToast(err.data?.error || err.message, "error");
  }
});

async function loadAdmin() {
  const { brands, categories } = await loadFilters();
  fillSelect('select[name="brand_id"]', brands);
  fillSelect('select[name="category_id"]', categories);
  renderList("#brands-list", brands, deleteBrand);
  renderList("#categories-list", categories, deleteCategory);
}

function fillSelect(sel, items) {
  const el = document.querySelector(sel);
  el.innerHTML = items.map((i) => `<option value="${i.id}">${escapeHtml(i.name)}</option>`).join("");
}

function renderList(sel, items, onDelete) {
  const ul = $(sel);
  ul.innerHTML = items
    .map(
      (i) => `
    <li>
      <span>${escapeHtml(i.name)}</span>
      <button class="btn btn--danger" data-id="${i.id}">Delete</button>
    </li>
  `
    )
    .join("");

  ul.querySelectorAll("button").forEach((btn) => {
    btn.addEventListener("click", () => onDelete(btn.dataset.id));
  });
}

async function deleteBrand(id) {
  await api(`/brands/${id}`, { method: "DELETE" });
  showToast("Brand deleted");
  loadAdmin();
}

async function deleteCategory(id) {
  await api(`/categories/${id}`, { method: "DELETE" });
  showToast("Category deleted");
  loadAdmin();
}

$("#brand-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const name = new FormData(e.target).get("name");
  await api("/brands", { method: "POST", body: JSON.stringify({ name }) });
  e.target.reset();
  showToast("Brand added");
  loadAdmin();
});

$("#category-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const name = new FormData(e.target).get("name");
  await api("/categories", { method: "POST", body: JSON.stringify({ name }) });
  e.target.reset();
  showToast("Category added");
  loadAdmin();
});

$("#product-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const fd = new FormData(e.target);
  await api("/products", {
    method: "POST",
    body: JSON.stringify({
      name: fd.get("name"),
      price: Number(fd.get("price")),
      stock: Number(fd.get("stock")),
      brand_id: Number(fd.get("brand_id")),
      category_id: Number(fd.get("category_id")),
    }),
  });
  e.target.reset();
  showToast("Product created");
  loadAdmin();
});

if (state.token) {
  updateAuthUI();
} else {
  showView("login");
  $("#nav-auth").classList.remove("hidden");
}
