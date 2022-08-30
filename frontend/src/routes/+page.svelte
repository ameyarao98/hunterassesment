<script lang="ts">
	import { SvelteToast, toast } from '@zerodevx/svelte-toast';
	import { GraphQLClient, gql } from 'graphql-request';

	let gameJoined: boolean;

	let jwt: string;

	let signupUsername: string;
	let signupPassword: string;

	let loginUsername: string;
	let loginPassword: string;

	let graphqlClient: GraphQLClient;
	let userResourceData: {
		resourceName: string;
		factoryLevel: number;
		amount: number;
		productionRate: number;
		timeUntilUpgradeComplete: number;
	}[];
	let factoryData: {
		iron: { factoryLevel: number; productionPerSecond: number; nextUpgradeDuration: number }[];
		copper: { factoryLevel: number; productionPerSecond: number; nextUpgradeDuration: number }[];
		gold: { factoryLevel: number; productionPerSecond: number; nextUpgradeDuration: number }[];
	};

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
			await login(signupUsername, signupPassword);
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
			toast.push(`Signed in as ${username}!`, {
				theme: {
					'--toastBackground': '#077023',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
			jwt = await response.text();
			await setupGame();
		}
	}

	async function setupGame() {
		graphqlClient = new GraphQLClient('http://localhost:3000/graphql', {
			headers: { creds: jwt }
		});
		await enterGame();
		await getFactoryData();
		gameJoined = true;
	}

	async function enterGame() {
		const query = gql`
			mutation {
				enterGame
			}
		`;
		try {
			await graphqlClient.request(query);
			toast.push("New factories created!", {
				theme: {
					'--toastBackground': '#077023',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
		} catch (error) {
			toast.push("Factories restored!", {
				theme: {
					'--toastBackground': '#077023',
					'--toastColor': 'white',
					'--toastBarBackground': 'white'
				}
			});
		}
	}

	async function getUserResourceData() {
		const query = gql`
			query {
				userResourceData {
					resourceName
					factoryLevel
					amount
					productionRate
					timeUntilUpgradeComplete
				}
			}
		`;
		userResourceData = (await graphqlClient.request(query)).userResourceData;
	}

	async function getFactoryData() {
		const query = gql`
			query {
				factoryData {
					iron {
						factoryLevel
						productionPerSecond
						nextUpgradeDuration
					}
					copper {
						factoryLevel
						productionPerSecond
						nextUpgradeDuration
					}
					gold {
						factoryLevel
						productionPerSecond
						nextUpgradeDuration
					}
				}
			}
		`;
		factoryData = (await graphqlClient.request(query)).factoryData;
	}
</script>

<h1>Welcome to the mining game!</h1>

{#if !gameJoined}
	<div class="login">
		<input class="field" bind:value={loginUsername} />
		<input class="field" type="password" bind:value={loginPassword} />
		<button class="first" on:click={() => login(loginUsername, loginPassword)}> Login </button>
	</div>

	<div class="signup">
		<input class="field" bind:value={signupUsername} />
		<input class="field" type="password" bind:value={signupPassword} />
		<button class="first" on:click={() => signup(signupUsername, signupPassword)}> Signup </button>
	</div>
{/if}

{#if gameJoined}
	<div class="factoryData">
		<table class="ironFactoryData">
			<tr>
				<th colspan="3">Iron</th>
			</tr>
			<tr>
				<th>Level</th>
				<th>Production Rate</th>
				<th>Next Upgrade Duration</th>
			</tr>
			{#each factoryData.iron as { factoryLevel, productionPerSecond, nextUpgradeDuration }}
				<tr>
					<th>{factoryLevel}</th>
					<th>{productionPerSecond}/s</th>
					<th>{nextUpgradeDuration} seconds</th>
				</tr>
			{/each}
		</table>
		<table class="copperFactoryData">
			<tr>
				<th colspan="3">Copper</th>
			</tr>
			<tr>
				<th>Level</th>
				<th>Production Rate</th>
				<th>Next Upgrade Duration</th>
			</tr>
			{#each factoryData.copper as { factoryLevel, productionPerSecond, nextUpgradeDuration }}
				<tr>
					<th>{factoryLevel}</th>
					<th>{productionPerSecond}/s</th>
					<th>{nextUpgradeDuration} seconds</th>
				</tr>
			{/each}
		</table>
		<table class="goldFactoryData">
			<tr>
				<th colspan="3">Gold</th>
			</tr>
			<tr>
				<th>Level</th>
				<th>Production Rate</th>
				<th>Next Upgrade Duration</th>
			</tr>
			{#each factoryData.gold as { factoryLevel, productionPerSecond, nextUpgradeDuration }}
				<tr>
					<th>{factoryLevel}</th>
					<th>{productionPerSecond}/s</th>
					<th>{nextUpgradeDuration} seconds</th>
				</tr>
			{/each}
		</table>
	</div>
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
		left: 1%;
	}

	.signup {
		display: grid;
		position: absolute;
		bottom: 2%;
		right: 1%;
	}

	.factoryData {
		width: 100%;
		display: flex;
		justify-content: center;
		bottom: 2%;
		position: absolute;
	}
	table {
		display: inline-table;
		margin-left: 5vw;
		margin-right: 5vw;
	}
	th,
	tr {
		border: 1px solid black;
	}

	.ironFactoryData {
		background-color: #a19d94;
	}
	.copperFactoryData {
		background-color: #b87333;
	}

	.goldFactoryData {
		background-color: #ffd700;
	}

	.jwt {
		font-size: 1vh;
		color: rgb(123, 255, 0);
		position: absolute;
		bottom: 0;
		text-align: center;
		opacity: 25%;
		z-index: -123123;
	}

	button {
		box-sizing: border-box;
		-webkit-appearance: none;
		-moz-appearance: none;
		appearance: none;
		background-color: transparent;
		border: 2px solid #ffd90f;
		border-radius: 0.6em;
		color: #ffd90f;
		cursor: pointer;
		display: -webkit-box;
		display: -webkit-flex;
		display: -ms-flexbox;
		display: flex;
		-webkit-align-self: center;
		-ms-flex-item-align: center;
		align-self: center;
		font-size: 1rem;
		font-weight: 400;
		line-height: 1;
		margin: 20px;
		padding: 1.2em 2.8em;
		text-decoration: none;
		text-align: center;
		text-transform: uppercase;
		font-family: 'Montserrat', sans-serif;
		font-weight: 700;
	}
	button:hover,
	button:focus {
		color: black;
		outline: 0;
	}

	.first {
		-webkit-transition: box-shadow 300ms ease-in-out, color 300ms ease-in-out;
		transition: box-shadow 300ms ease-in-out, color 300ms ease-in-out;
	}
	.first:hover {
		box-shadow: 0 0 40px 40px #ffd90f inset;
	}

	input {
		padding: 0.6em 0em;
		text-align: center;
	}
</style>
