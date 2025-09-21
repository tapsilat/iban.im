<template>
  <div class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow">
    <div id="scan">
  <ul class="space-y-4">
        
        <li style="--color: var(--primary-2)">
          <div id="preview">
            <img id="iban_img" style="max-height: 40px;" alt="" v-if="url" :src="url" />
          </div>
        </li>
        <li>
          <input type="file" @change="onFileChange" />
        </li>
        <template v-if="scanning">
          <div class="flex justify-center py-4">
            <svg class="animate-spin h-6 w-6 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
            </svg>
          </div>
        </template>
        <template v-else>
          <li>
            <button v-on:click="recognize" class="px-4 py-2 rounded bg-blue-600 text-white">recognize</button>
          </li>
          <li style="--color: var(--primary-3)">
            Text in the image:
            <div id="iban_raw">{{ibanRaw}}</div>
          </li>
          <li style="--color: var(--primary-3)">
            possible IBAN:
            <strong id="iban">{{iban}}</strong>
          </li>
        </template>

      </ul>
    </div>
  </div>
</template>

<script>
/* eslint-disable */
import { createWorker, PSM, OEM } from "tesseract.js";
const worker = createWorker({
  logger: (m) => console.log(m),
});

export default {
  name: "Scan",
  data () {
    return {
      url: null,
      text: null,
      iban: null,
      ibanRaw : null,
      scanning : false,
    };
  },
  methods: {
    async recognize() {
      this.scanning = true;
      const iban_img = document.getElementById("iban_img");
      await worker.load();
      await worker.loadLanguage("eng");
      await worker.initialize("eng", OEM.LSTM_ONLY);
      await worker.setParameters({
        tessedit_pageseg_mode: PSM.SINGLE_BLOCK,
      });
      let {
        data: { text },
      } = await worker.recognize(iban_img);
      this.scanning = false;
      console.log(text);
      this.ibanRaw = text;
      text = text.replace(/ /g, '')
      text = text.match(
        /[a-zA-Z]{2}[0-9]{2}[a-zA-Z0-9]{4}[0-9]{7}([a-zA-Z0-9]?){0,16}/m
      );
      if(text.length > 0) {
        this.iban = text[0];
      }
    },
    onFileChange(e) {
      const file = e.target.files[0];
      this.url = URL.createObjectURL(file);
    },
  },
};
</script>