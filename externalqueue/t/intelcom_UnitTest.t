#!/usr/bin/env ferite
//#!/bin/sh\nferite externalqueue.fe
uses 'externalqueue';
uses 'intelecom.feh';
uses 'console';
uses 'date';
uses 'workflow';
uses 'tap.feh';
uses '/cention/src/cention/cention-services/scripts.feh';

//unit test cases for intelcom interface
//this will only work if some form of soap server
//is available for testing


function authenticate_test_01(){
	boolean isOk;
	array options = [ 'server-address' => 'http://localhost:12347',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.authenticate();

	is(IntelecomQueue.tokenString, "blahblahblahtoken", "Authenticated");
	is(isOk, true,"authenticate_test_01 return true");
}

function addErrand_test_01() {
	boolean isOk;
	object errand = new Workflow.Errand();
	object user = new Workflow.User();
	array options = [ 'server-address' => 'http://localhost:12347',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'trigger-uri' =>'http://localhost:12347/action',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 2;
	errand.externalID = 0;
	errand.mail = new Workflow.Mail();
	errand.mail.from = new Workflow.MailOrigin();
	errand.mail.from.name = "finn";
	errand.mail.from.emailAddress = "finn@adventuretime.com";
	errand.mail.subject = "candy kingdom";
	errand.timestampArrive = Date.now();
	user.timezoneID = 0;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.addErrand(errand, user, 0);
	is(isOk, true,"addErrand_test_01 return true");
	is(errand.externalID, 4006, "addErrand_test_01 returned correct externalID");
}

function delErrand_test_01() {
	boolean isOk;
	object errand = new Workflow.Errand();
	array options = [ 'server-address' => 'http://localhost:12347',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 3;
	errand.externalID = 4006;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.removeErrand(errand, 0, false);
	is(isOk, true,"delErrand_test_01 return true");
	is(errand.externalID, 0, "delErrand_test_01 returned correct externalID");
}

function pullErrand_test_01() {
	boolean isOk;
	object errand = new Workflow.Errand();
	array options = [ 'server-address' => 'http://localhost:12347',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 3;
	errand.externalID = 4006;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.pullErrand(errand);
	is(isOk, true,"getErrand_test_01 return true");
}

function authenticate_realtest_01(){
	boolean isOk;
	array options = [ 'server-address' => 'https://api.intele.com/Connect/ContactCentre/Agent.svc',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.authenticate();

	is((IntelecomQueue.tokenString.length() > 0), 
		true, "Authenticated");
	is(isOk, true,"authenticate_realtest_01 return true");
}

function addErrand_realtest_01() {
	boolean isOk;
	object errand = new Workflow.Errand();
	object user = new Workflow.User();
	array options = [ 'server-address' => 'https://api.intele.com/Connect/ContactCentre/Agent.svc',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 5783;
	errand.externalID = 0;
	errand.mail = new Workflow.Mail();
	errand.mail.from = new Workflow.MailOrigin();
	errand.mail.from.name = "finn";
	errand.mail.from.emailAddress = "finn@adventuretime.com";
	errand.mail.subject = "candy kingdom";
	errand.timestampArrive = Date.now();
	user.timezoneID = 0;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.addErrand(errand, user, 0);
	Console.println("external Id:"+errand.externalID);
	is(isOk, true,"addErrand_realtest_01 return true");
	is((errand.externalID > 0), true, 
		"addErrand_realtest_01 returned correct externalID");
}

function delErrand_realtest_01() {
	boolean isOk;
	object errand = new Workflow.Errand();
	array options = [ 'server-address' => 'https://api.intele.com/Connect/ContactCentre/Agent.svc',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 3;
	errand.externalID = 601074221;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.removeErrand(errand, 0, false);
	Console.println("external Id:"+errand.externalID);
	is(isOk, true,"delErrand_realtest_01 return true");
	is((errand.externalID > 0), false, 
		"delErrand_realtest_01 returned correct externalID");
}

function pullErrand_realtest_01() {
	boolean isOk = false;
	object errand = new Workflow.Errand();
	array options = [ 'server-address' => 'https://api.intele.com/Connect/ContactCentre/Agent.svc',
		'soap-action-url' => 'ContactCentreWebServices/IAgent',
		'open-errand-url' =>'',
		'customerkey'=> '99218',
		'username' => 'Cention1',
		'password'=> 'Tobias1',
		'access-point'=> 'test.intelecom.1@cention.se'
	];
	errand.id = 3;
	errand.externalID = 601977154;
	IntelecomQueue.configure(options);
	isOk = IntelecomQueue.pullErrand(errand);
	is(isOk, true,"getErrand_test_01 return true");
}

//test starts here
//authenticate_test_01();
//addErrand_test_01();
//delErrand_test_01();
//pullErrand_test_01();

//test connecting to intelcom server
//Caution. too many authentication failure will result in account being locked
//authenticate_realtest_01();
//addErrand_realtest_01();
//delErrand_realtest_01();
//pullErrand_realtest_01();
