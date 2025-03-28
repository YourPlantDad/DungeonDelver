<script>
    import { onMount } from 'svelte';
  
    let message = "Laden...";
  
    onMount(async () => {
      try {
        const res = await fetch("http://localhost:8080/api/hello");
        const data = await res.json();
        message = data.message;
      } catch (error) {
        console.error("Fout bij het ophalen van bericht:", error);
        message = "Er is iets misgegaan";
      }
    });
  </script>
  
  <main>
    <h1>{message}</h1>
  </main>
  