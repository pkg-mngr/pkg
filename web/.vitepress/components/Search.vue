<script setup lang="ts">
import { ref, onMounted } from "vue";

const value = ref("");
let packages = [];
let listItems: NodeListOf<HTMLLIElement>;

onMounted(() => {
    listItems = document.querySelectorAll<HTMLLIElement>("li[data-name]");
    listItems.forEach((item) => {
        const { name, desc } = item.dataset;
        packages.push(`${name.toLowerCase()} ${desc.toLowerCase()}`.trim());
    });
});

function handleInput() {
    if (value.value.trim() === "") {
        listItems.forEach((item) => (item.style.display = "list-item"));
        return;
    }
    listItems.forEach((item, index) => {
        if (!packages[index].includes(value.value.trim().toLowerCase())) {
            item.style.display = "none";
        } else {
            item.style.display = "list-item";
        }
    });
}
</script>

<template>
    <input
        type="search"
        name="Package search"
        id="package-search"
        placeholder="Search for packages..."
        v-model="value"
        @input="handleInput"
    />
</template>

<style scoped>
input {
    background-color: var(--vp-c-bg-alt);
    border-radius: 0.5rem;
    width: 100%;
    height: 2.5rem;
    padding: 0.75rem;
    font-size: 1rem;
    margin: 1rem 0;

    &::placeholder {
        font-size: 0.875rem;
        font-weight: 500;
        color: var(--vp-c-text-2);
    }
}
</style>
