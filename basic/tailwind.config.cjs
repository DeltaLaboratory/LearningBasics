import { join } from 'path'
import forms from '@tailwindcss/forms'
import typography from '@tailwindcss/typography'
import skeleton from '@skeletonlabs/skeleton/tailwind/skeleton.cjs'

/** @type {import('tailwindcss').Config} */
module.exports = {
	darkMode: 'class',
	content: ['./src/**/*.{html,js,svelte,ts}', join(require.resolve('@skeletonlabs/skeleton'), '../**/*.{html,js,svelte,ts}')],
	theme: {
		extend: {},
	},
	plugins: [forms,typography,...skeleton()],
	purge: {
		options: {
			safelist: {
				standard: [
					'text-2xl',
					'text-3xl',
					'text-4xl',
					'text-5xl',
					'text-6xl',
					'sm:text-2xl',
					'sm:text-3xl',
					'sm:text-4xl',
					'sm:text-5xl',
					'sm:text-6xl',
					'md:text-2xl',
					'md:text-3xl',
					'md:text-4xl',
					'md:text-5xl',
					'md:text-6xl',
					'lg:text-2xl',
					'lg:text-3xl',
					'lg:text-4xl',
					'lg:text-5xl',
					'lg:text-6xl',
				],
			},
		}
	}
}
