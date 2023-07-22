export class AppHealthcheck {
	enabled = false;
	health = 'unknown'; // TODO: make this an enum ?
	interval = 0;
	timeout = 0;
}

export class App {
	name = '';
	link = '';
	group = '';
	description = '';
	icon = '';
	healthcheck = new AppHealthcheck();
}

export class AppGroup {
	name = '';
	apps: App[] = [];
}

export class AppSettings {
	name = 'simplydash';
}
