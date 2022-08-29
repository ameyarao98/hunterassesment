<script lang="ts">
	import { SvelteToast, toast } from '@zerodevx/svelte-toast';

	let gameJoined: boolean;

	let jwt: string;

	let signupUsername: string;
	let signupPassword: string;

	let loginUsername: string;
	let loginPassword: string;

	async function signup(username: string, password: string) {
		const response = await fetch('http://localhost:8000/signup', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ username: username, password: password })
		});
		if (response.status != 204) {
			toast.push('Signup failed!', {
				theme: {
					'--toastBackground': '#8b0000',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
		} else {
			toast.push('Signup succeeded!', {
				theme: {
					'--toastBackground': '#077023',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
			signupUsername = '';
			signupPassword = '';
		}
	}
	async function login(username: string, password: string) {
		const response = await fetch('http://localhost:8000/auth', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ username: username, password: password })
		});
		if (response.status != 200) {
			toast.push('Login failed!', {
				theme: {
					'--toastBackground': '#8b0000',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
		} else {
			jwt = await response.text();
			toast.push(`Signed in as ${loginUsername}!`, {
				theme: {
					'--toastBackground': '#077023',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
			loginUsername = '';
			loginPassword = '';
			gameJoined = true;
		}
	}
</script>

<h1>Welcome to the mining game!</h1>

{#if !gameJoined}
	<div class="login">
		<input class="field" maxlength="20" bind:value={loginUsername} />
		<input class="field" maxlength="20" bind:value={loginPassword} />
		<button on:click={() => login(loginUsername, loginPassword)}> Login </button>
	</div>
{/if}

{#if !gameJoined}
	<div class="signup">
		<input class="field" maxlength="20" bind:value={signupUsername} />
		<input class="field" maxlength="20" bind:value={signupPassword} />
		<button on:click={() => signup(signupUsername, signupPassword)}> Signup </button>
	</div>
{/if}

{#if gameJoined}
	<p class="jwt">{jwt}</p>
{/if}

<SvelteToast />

<style>
	:global(body) {
		background-image: url('mining.gif');
		background-size: cover;
		height: 100vh;
		padding: 0;
		margin: 0;
	}
	h1 {
		margin-top: 0;
		-webkit-text-stroke: 0.3px black;
		font-size: 10vh;
		text-align: center;
		font-family: 'Silkscreen', sans-serif;
		font-weight: bold;
		color: #ffd90f;
	}

	.login {
		display: grid;
		position: absolute;
		bottom: 2%;
		left: 0;
	}

	.signup {
		display: grid;
		position: absolute;
		bottom: 2%;
		right: 0;
	}

	.jwt {
		font-size: 1vh;
		color: rgb(123, 255, 0);
		position: absolute;
		bottom: 0;
		text-align: center;
	}
</style>
