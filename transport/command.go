package transport

import (
	"fmt"
)

type Command uint16

// Provides information on what triggered the error.
const GW_ERROR_NTF Command = 0x0000

// Request gateway to reboot.
const GW_REBOOT_REQ Command = 0x0001

// Acknowledge to GW_REBOOT_REQ command.
const GW_REBOOT_CFM Command = 0x0002

// Request gateway to clear system table, scene table and set Ethernet settings to factory default. Gateway will reboot.
const GW_SET_FACTORY_DEFAULT_REQ Command = 0x0003

// Acknowledge to GW_SET_FACTORY_DEFAULT_REQ command.
const GW_SET_FACTORY_DEFAULT_CFM Command = 0x0004

// Request version information.
const GW_GET_VERSION_REQ Command = 0x0008

// Acknowledge to GW_GET_VERSION_REQ command.
const GW_GET_VERSION_CFM Command = 0x0009

// Request KLF 200 API protocol version.
const GW_GET_PROTOCOL_VERSION_REQ Command = 0x000A

// Acknowledge to GW_GET_PROTOCOL_VERSION_REQ command.
const GW_GET_PROTOCOL_VERSION_CFM Command = 0x000B

// Request the state of the gateway
const GW_GET_STATE_REQ Command = 0x000C

// Acknowledge to GW_GET_STATE_REQ command.
const GW_GET_STATE_CFM Command = 0x000D

// Request gateway to leave learn state.
const GW_LEAVE_LEARN_STATE_REQ Command = 0x000E

// Acknowledge to GW_LEAVE_LEARN_STATE_REQ command.
const GW_LEAVE_LEARN_STATE_CFM Command = 0x000F

// Request network parameters.
const GW_GET_NETWORK_SETUP_REQ Command = 0x00E0

// Acknowledge to GW_GET_NETWORK_SETUP_REQ.
const GW_GET_NETWORK_SETUP_CFM Command = 0x00E1

// Set network parameters.
const GW_SET_NETWORK_SETUP_REQ Command = 0x00E2

// Acknowledge to GW_SET_NETWORK_SETUP_REQ.
const GW_SET_NETWORK_SETUP_CFM Command = 0x00E3

// Request a list of nodes in the gateways system table.
const GW_CS_GET_SYSTEMTABLE_DATA_REQ Command = 0x0100

// Acknowledge to GW_CS_GET_SYSTEMTABLE_DATA_REQ
const GW_CS_GET_SYSTEMTABLE_DATA_CFM Command = 0x0101

// Acknowledge to GW_CS_GET_SYSTEM_TABLE_DATA_REQList of nodes in the gateways systemtable.
const GW_CS_GET_SYSTEMTABLE_DATA_NTF Command = 0x0102

// Start CS DiscoverNodes macro in KLF200.
const GW_CS_DISCOVER_NODES_REQ Command = 0x0103

// Acknowledge to GW_CS_DISCOVER_NODES_REQ command.
const GW_CS_DISCOVER_NODES_CFM Command = 0x0104

// Acknowledge to GW_CS_DISCOVER_NODES_REQ command.
const GW_CS_DISCOVER_NODES_NTF Command = 0x0105

// Remove one or more nodes in the systemtable.
const GW_CS_REMOVE_NODES_REQ Command = 0x0106

// Acknowledge to GW_CS_REMOVE_NODES_REQ.
const GW_CS_REMOVE_NODES_CFM Command = 0x0107

// Clear systemtable and delete system key.
const GW_CS_VIRGIN_STATE_REQ Command = 0x0108

// Acknowledge to GW_CS_VIRGIN_STATE_REQ.
const GW_CS_VIRGIN_STATE_CFM Command = 0x0109

// Setup KLF200 to get or give a system to or from another io-homecontrol® remote control. By a system means all nodes in the systemtable and the system key.
const GW_CS_CONTROLLER_COPY_REQ Command = 0x010A

// Acknowledge to GW_CS_CONTROLLER_COPY_REQ.
const GW_CS_CONTROLLER_COPY_CFM Command = 0x010B

// Acknowledge to GW_CS_CONTROLLER_COPY_REQ.
const GW_CS_CONTROLLER_COPY_NTF Command = 0x010C

// Cancellation of system copy to other controllers.
const GW_CS_CONTROLLER_COPY_CANCEL_NTF Command = 0x010D

// Receive system key from another controller.
const GW_CS_RECEIVE_KEY_REQ Command = 0x010E

// Acknowledge to GW_CS_RECEIVE_KEY_REQ.
const GW_CS_RECEIVE_KEY_CFM Command = 0x010F

// Acknowledge to GW_CS_RECEIVE_KEY_REQ with status.
const GW_CS_RECEIVE_KEY_NTF Command = 0x0110

// Information on Product Generic Configuration job initiated by press on PGC button.
const GW_CS_PGC_JOB_NTF Command = 0x0111

// Broadcasted to all clients and gives information about added and removed actuator nodes in system table.
const GW_CS_SYSTEM_TABLE_UPDATE_NTF Command = 0x0112

// Generate new system key and update actuators in systemtable.
const GW_CS_GENERATE_NEW_KEY_REQ Command = 0x0113

// Acknowledge to GW_CS_GENERATE_NEW_KEY_REQ.
const GW_CS_GENERATE_NEW_KEY_CFM Command = 0x0114

// Acknowledge to GW_CS_GENERATE_NEW_KEY_REQ with status.
const GW_CS_GENERATE_NEW_KEY_NTF Command = 0x0115

// Update key in actuators holding an old key.
const GW_CS_REPAIR_KEY_REQ Command = 0x0116

// Acknowledge to GW_CS_REPAIR_KEY_REQ.
const GW_CS_REPAIR_KEY_CFM Command = 0x0117

// Acknowledge to GW_CS_REPAIR_KEY_REQ with status.
const GW_CS_REPAIR_KEY_NTF Command = 0x0118

// Request one or more actuator to open for configuration.
const GW_CS_ACTIVATE_CONFIGURATION_MODE_REQ Command = 0x0119

// Acknowledge to GW_CS_ACTIVATE_CONFIGURATION_MODE_REQ.
const GW_CS_ACTIVATE_CONFIGURATION_MODE_CFM Command = 0x011A

// Request extended information of one specific actuator node.
const GW_GET_NODE_INFORMATION_REQ Command = 0x0200

// Acknowledge to GW_GET_NODE_INFORMATION_REQ.
const GW_GET_NODE_INFORMATION_CFM Command = 0x0201

// Acknowledge to GW_GET_NODE_INFORMATION_REQ.
const GW_GET_NODE_INFORMATION_NTF Command = 0x0210

// Request extended information of all nodes.
const GW_GET_ALL_NODES_INFORMATION_REQ Command = 0x0202

// Acknowledge to GW_GET_ALL_NODES_INFORMATION_REQ
const GW_GET_ALL_NODES_INFORMATION_CFM Command = 0x0203

// Acknowledge to GW_GET_ALL_NODES_INFORMATION_REQ. Holds node information
const GW_GET_ALL_NODES_INFORMATION_NTF Command = 0x0204

// Acknowledge to GW_GET_ALL_NODES_INFORMATION_REQ. No more nodes.
const GW_GET_ALL_NODES_INFORMATION_FINISHED_NTF Command = 0x0205

// Set node variation.
const GW_SET_NODE_VARIATION_REQ Command = 0x0206

// Acknowledge to GW_SET_NODE_VARIATION_REQ.
const GW_SET_NODE_VARIATION_CFM Command = 0x0207

// Set node name.
const GW_SET_NODE_NAME_REQ Command = 0x0208

// Acknowledge to GW_SET_NODE_NAME_REQ.
const GW_SET_NODE_NAME_CFM Command = 0x0209

// Information has been updated.
const GW_NODE_INFORMATION_CHANGED_NTF Command = 0x020C

// Information has been updated.
const GW_NODE_STATE_POSITION_CHANGED_NTF Command = 0x0211

// Set search order and room placement.
const GW_SET_NODE_ORDER_AND_PLACEMENT_REQ Command = 0x020D

// Acknowledge to GW_SET_NODE_ORDER_AND_PLACEMENT_REQ.
const GW_SET_NODE_ORDER_AND_PLACEMENT_CFM Command = 0x020E

// Request information about all defined groups.
const GW_GET_GROUP_INFORMATION_REQ Command = 0x0220

// Acknowledge to GW_GET_GROUP_INFORMATION_REQ.
const GW_GET_GROUP_INFORMATION_CFM Command = 0x0221

// Acknowledge to GW_GET_NODE_INFORMATION_REQ.
const GW_GET_GROUP_INFORMATION_NTF Command = 0x0230

// Change an existing group.
const GW_SET_GROUP_INFORMATION_REQ Command = 0x0222

// Acknowledge to GW_SET_GROUP_INFORMATION_REQ.
const GW_SET_GROUP_INFORMATION_CFM Command = 0x0223

// Broadcast to all, about group information of a group has been changed.
const GW_GROUP_INFORMATION_CHANGED_NTF Command = 0x0224

// Delete a group.
const GW_DELETE_GROUP_REQ Command = 0x0225

// Acknowledge to GW_DELETE_GROUP_INFORMATION_REQ.
const GW_DELETE_GROUP_CFM Command = 0x0226

// Request new group to be created.
const GW_NEW_GROUP_REQ Command = 0x0227

const GW_NEW_GROUP_CFM Command = 0x0228

// Request information about all defined groups.
const GW_GET_ALL_GROUPS_INFORMATION_REQ Command = 0x0229

// Acknowledge to GW_GET_ALL_GROUPS_INFORMATION_REQ.
const GW_GET_ALL_GROUPS_INFORMATION_CFM Command = 0x022A

// Acknowledge to GW_GET_ALL_GROUPS_INFORMATION_REQ.
const GW_GET_ALL_GROUPS_INFORMATION_NTF Command = 0x022B

// Acknowledge to GW_GET_ALL_GROUPS_INFORMATION_REQ.
const GW_GET_ALL_GROUPS_INFORMATION_FINISHED_NTF Command = 0x022C

// GW_GROUP_DELETED_NTF is broadcasted to all, when a group has been removed.
const GW_GROUP_DELETED_NTF Command = 0x022D

// Enable house status monitor.
const GW_HOUSE_STATUS_MONITOR_ENABLE_REQ Command = 0x0240

// Acknowledge to GW_HOUSE_STATUS_MONITOR_ENABLE_REQ.
const GW_HOUSE_STATUS_MONITOR_ENABLE_CFM Command = 0x0241

// Disable house status monitor.
const GW_HOUSE_STATUS_MONITOR_DISABLE_REQ Command = 0x0242

// Acknowledge to GW_HOUSE_STATUS_MONITOR_DISABLE_REQ.
const GW_HOUSE_STATUS_MONITOR_DISABLE_CFM Command = 0x0243

// Send activating command direct to one or more io-homecontrol® nodes.
const GW_COMMAND_SEND_REQ Command = 0x0300

// Acknowledge to GW_COMMAND_SEND_REQ.
const GW_COMMAND_SEND_CFM Command = 0x0301

// Gives run status for io-homecontrol® node.
const GW_COMMAND_RUN_STATUS_NTF Command = 0x0302

// Gives remaining time before io-homecontrol® node enter target position.
const GW_COMMAND_REMAINING_TIME_NTF Command = 0x0303

// Command send, Status request, Wink, Mode or Stop session is finished.
const GW_SESSION_FINISHED_NTF Command = 0x0304

// Get status request from one or more io-homecontrol® nodes.
const GW_STATUS_REQUEST_REQ Command = 0x0305

// Acknowledge to GW_STATUS_REQUEST_REQ.
const GW_STATUS_REQUEST_CFM Command = 0x0306

// Acknowledge to GW_STATUS_REQUEST_REQ. Status request from one or more io-homecontrol® nodes.
const GW_STATUS_REQUEST_NTF Command = 0x0307

// Request from one or more io-homecontrol® nodes to Wink.
const GW_WINK_SEND_REQ Command = 0x0308

// Acknowledge to GW_WINK_SEND_REQ
const GW_WINK_SEND_CFM Command = 0x0309

// Status info for performed wink request.
const GW_WINK_SEND_NTF Command = 0x030A

// Set a parameter limitation in an actuator.
const GW_SET_LIMITATION_REQ Command = 0x0310

// Acknowledge to GW_SET_LIMITATION_REQ.
const GW_SET_LIMITATION_CFM Command = 0x0311

// Get parameter limitation in an actuator.
const GW_GET_LIMITATION_STATUS_REQ Command = 0x0312

// Acknowledge to GW_GET_LIMITATION_STATUS_REQ.
const GW_GET_LIMITATION_STATUS_CFM Command = 0x0313

// Hold information about limitation.
const GW_LIMITATION_STATUS_NTF Command = 0x0314

// Send Activate Mode to one or more io-homecontrol® nodes.
const GW_MODE_SEND_REQ Command = 0x0320

// Acknowledge to GW_MODE_SEND_REQ
const GW_MODE_SEND_CFM Command = 0x0321

// Notify with Mode activation info.
const GW_MODE_SEND_NTF Command = 0x0322

// Prepare gateway to record a scene.
const GW_INITIALIZE_SCENE_REQ Command = 0x0400

// Acknowledge to GW_INITIALIZE_SCENE_REQ.
const GW_INITIALIZE_SCENE_CFM Command = 0x0401

// Acknowledge to GW_INITIALIZE_SCENE_REQ.
const GW_INITIALIZE_SCENE_NTF Command = 0x0402

// Cancel record scene process.
const GW_INITIALIZE_SCENE_CANCEL_REQ Command = 0x0403

// Acknowledge to GW_INITIALIZE_SCENE_CANCEL_REQ command.
const GW_INITIALIZE_SCENE_CANCEL_CFM Command = 0x0404

// Store actuator positions changes since GW_INITIALIZE_SCENE, as a scene.
const GW_RECORD_SCENE_REQ Command = 0x0405

// Acknowledge to GW_RECORD_SCENE_REQ.
const GW_RECORD_SCENE_CFM Command = 0x0406

// Acknowledge to GW_RECORD_SCENE_REQ.
const GW_RECORD_SCENE_NTF Command = 0x0407

// Delete a recorded scene.
const GW_DELETE_SCENE_REQ Command = 0x0408

// Acknowledge to GW_DELETE_SCENE_REQ.
const GW_DELETE_SCENE_CFM Command = 0x0409

// Request a scene to be renamed.
const GW_RENAME_SCENE_REQ Command = 0x040A

// Acknowledge to GW_RENAME_SCENE_REQ.
const GW_RENAME_SCENE_CFM Command = 0x040B

// Request a list of scenes.
const GW_GET_SCENE_LIST_REQ Command = 0x040C

// Acknowledge to GW_GET_SCENE_LIST.
const GW_GET_SCENE_LIST_CFM Command = 0x040D

// Acknowledge to GW_GET_SCENE_LIST.
const GW_GET_SCENE_LIST_NTF Command = 0x040E

// Request extended information for one given scene.
const GW_GET_SCENE_INFOAMATION_REQ Command = 0x040F

// Acknowledge to GW_GET_SCENE_INFOAMATION_REQ.
const GW_GET_SCENE_INFOAMATION_CFM Command = 0x0410

// Acknowledge to GW_GET_SCENE_INFOAMATION_REQ.
const GW_GET_SCENE_INFOAMATION_NTF Command = 0x0411

// Request gateway to enter a scene.
const GW_ACTIVATE_SCENE_REQ Command = 0x0412

// Acknowledge to GW_ACTIVATE_SCENE_REQ.
const GW_ACTIVATE_SCENE_CFM Command = 0x0413

// Request all nodes in a given scene to stop at their current position.
const GW_STOP_SCENE_REQ Command = 0x0415

// Acknowledge to GW_STOP_SCENE_REQ.
const GW_STOP_SCENE_CFM Command = 0x0416

// A scene has either been changed or removed.
const GW_SCENE_INFORMATION_CHANGED_NTF Command = 0x0419

// Activate a product group in a given direction.
const GW_ACTIVATE_PRODUCTGROUP_REQ Command = 0x0447

// Acknowledge to GW_ACTIVATE_PRODUCTGROUP_REQ.
const GW_ACTIVATE_PRODUCTGROUP_CFM Command = 0x0448

// Acknowledge to GW_ACTIVATE_PRODUCTGROUP_REQ.
const GW_ACTIVATE_PRODUCTGROUP_NTF Command = 0x0449

// Get list of assignments to all Contact Input to scene or product group.
const GW_GET_CONTACT_INPUT_LINK_LIST_REQ Command = 0x0460

// Acknowledge to GW_GET_CONTACT_INPUT_LINK_LIST_REQ.
const GW_GET_CONTACT_INPUT_LINK_LIST_CFM Command = 0x0461

// Set a link from a Contact Input to a scene or product group.
const GW_SET_CONTACT_INPUT_LINK_REQ Command = 0x0462

// Acknowledge to GW_SET_CONTACT_INPUT_LINK_REQ.
const GW_SET_CONTACT_INPUT_LINK_CFM Command = 0x0463

// Remove a link from a Contact Input to a scene.
const GW_REMOVE_CONTACT_INPUT_LINK_REQ Command = 0x0464

// Acknowledge to GW_REMOVE_CONTACT_INPUT_LINK_REQ.
const GW_REMOVE_CONTACT_INPUT_LINK_CFM Command = 0x0465

// Request header from activation log.
const GW_GET_ACTIVATION_LOG_HEADER_REQ Command = 0x0500

// Confirm header from activation log.
const GW_GET_ACTIVATION_LOG_HEADER_CFM Command = 0x0501

// Request clear all data in activation log.
const GW_CLEAR_ACTIVATION_LOG_REQ Command = 0x0502

// Confirm clear all data in activation log.
const GW_CLEAR_ACTIVATION_LOG_CFM Command = 0x0503

// Request line from activation log.
const GW_GET_ACTIVATION_LOG_LINE_REQ Command = 0x0504

// Confirm line from activation log.
const GW_GET_ACTIVATION_LOG_LINE_CFM Command = 0x0505

// Confirm line from activation log.
const GW_ACTIVATION_LOG_UPDATED_NTF Command = 0x0506

// Request lines from activation log.
const GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_REQ Command = 0x0507

// Error log data from activation log.
const GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_NTF Command = 0x0508

// Confirm lines from activation log.
const GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_CFM Command = 0x0509

// Request to set UTC time.
const GW_SET_UTC_REQ Command = 0x2000

// Acknowledge to GW_SET_UTC_REQ.
const GW_SET_UTC_CFM Command = 0x2001

// Set time zone and daylight savings rules.
const GW_RTC_SET_TIME_ZONE_REQ Command = 0x2002

// Acknowledge to GW_RTC_SET_TIME_ZONE_REQ.
const GW_RTC_SET_TIME_ZONE_CFM Command = 0x2003

// Request the local time based on current time zone and daylight savings rules.
const GW_GET_LOCAL_TIME_REQ Command = 0x2004

// Acknowledge to GW_RTC_SET_TIME_ZONE_REQ.
const GW_GET_LOCAL_TIME_CFM Command = 0x2005

// Enter password to authenticate request
const GW_PASSWORD_ENTER_REQ Command = 0x3000

// Acknowledge to GW_PASSWORD_ENTER_REQ
const GW_PASSWORD_ENTER_CFM Command = 0x3001

// Request password change.
const GW_PASSWORD_CHANGE_REQ Command = 0x3002

// Acknowledge to GW_PASSWORD_CHANGE_REQ.
const GW_PASSWORD_CHANGE_CFM Command = 0x3003

// Acknowledge to GW_PASSWORD_CHANGE_REQ. Broadcasted to all connected clients.
const GW_PASSWORD_CHANGE_NTF Command = 0x3004

func (cmd Command) String() string {
	switch cmd {
	case GW_ERROR_NTF:
		return "GW_ERROR_NTF"
	case GW_REBOOT_REQ:
		return "GW_REBOOT_REQ"
	case GW_REBOOT_CFM:
		return "GW_REBOOT_CFM"
	case GW_SET_FACTORY_DEFAULT_REQ:
		return "GW_SET_FACTORY_DEFAULT_REQ"
	case GW_SET_FACTORY_DEFAULT_CFM:
		return "GW_SET_FACTORY_DEFAULT_CFM"
	case GW_GET_VERSION_REQ:
		return "GW_GET_VERSION_REQ"
	case GW_GET_VERSION_CFM:
		return "GW_GET_VERSION_CFM"
	case GW_GET_PROTOCOL_VERSION_REQ:
		return "GW_GET_PROTOCOL_VERSION_REQ"
	case GW_GET_PROTOCOL_VERSION_CFM:
		return "GW_GET_PROTOCOL_VERSION_CFM"
	case GW_GET_STATE_REQ:
		return "GW_GET_STATE_REQ"
	case GW_GET_STATE_CFM:
		return "GW_GET_STATE_CFM"
	case GW_LEAVE_LEARN_STATE_REQ:
		return "GW_LEAVE_LEARN_STATE_REQ"
	case GW_LEAVE_LEARN_STATE_CFM:
		return "GW_LEAVE_LEARN_STATE_CFM"
	case GW_GET_NETWORK_SETUP_REQ:
		return "GW_GET_NETWORK_SETUP_REQ"
	case GW_GET_NETWORK_SETUP_CFM:
		return "GW_GET_NETWORK_SETUP_CFM"
	case GW_SET_NETWORK_SETUP_REQ:
		return "GW_SET_NETWORK_SETUP_REQ"
	case GW_SET_NETWORK_SETUP_CFM:
		return "GW_SET_NETWORK_SETUP_CFM"
	case GW_CS_GET_SYSTEMTABLE_DATA_REQ:
		return "GW_CS_GET_SYSTEMTABLE_DATA_REQ"
	case GW_CS_GET_SYSTEMTABLE_DATA_CFM:
		return "GW_CS_GET_SYSTEMTABLE_DATA_CFM"
	case GW_CS_GET_SYSTEMTABLE_DATA_NTF:
		return "GW_CS_GET_SYSTEMTABLE_DATA_NTF"
	case GW_CS_DISCOVER_NODES_REQ:
		return "GW_CS_DISCOVER_NODES_REQ"
	case GW_CS_DISCOVER_NODES_CFM:
		return "GW_CS_DISCOVER_NODES_CFM"
	case GW_CS_DISCOVER_NODES_NTF:
		return "GW_CS_DISCOVER_NODES_NTF"
	case GW_CS_REMOVE_NODES_REQ:
		return "GW_CS_REMOVE_NODES_REQ"
	case GW_CS_REMOVE_NODES_CFM:
		return "GW_CS_REMOVE_NODES_CFM"
	case GW_CS_VIRGIN_STATE_REQ:
		return "GW_CS_VIRGIN_STATE_REQ"
	case GW_CS_VIRGIN_STATE_CFM:
		return "GW_CS_VIRGIN_STATE_CFM"
	case GW_CS_CONTROLLER_COPY_REQ:
		return "GW_CS_CONTROLLER_COPY_REQ"
	case GW_CS_CONTROLLER_COPY_CFM:
		return "GW_CS_CONTROLLER_COPY_CFM"
	case GW_CS_CONTROLLER_COPY_NTF:
		return "GW_CS_CONTROLLER_COPY_NTF"
	case GW_CS_CONTROLLER_COPY_CANCEL_NTF:
		return "GW_CS_CONTROLLER_COPY_CANCEL_NTF"
	case GW_CS_RECEIVE_KEY_REQ:
		return "GW_CS_RECEIVE_KEY_REQ"
	case GW_CS_RECEIVE_KEY_CFM:
		return "GW_CS_RECEIVE_KEY_CFM"
	case GW_CS_RECEIVE_KEY_NTF:
		return "GW_CS_RECEIVE_KEY_NTF"
	case GW_CS_PGC_JOB_NTF:
		return "GW_CS_PGC_JOB_NTF"
	case GW_CS_SYSTEM_TABLE_UPDATE_NTF:
		return "GW_CS_SYSTEM_TABLE_UPDATE_NTF"
	case GW_CS_GENERATE_NEW_KEY_REQ:
		return "GW_CS_GENERATE_NEW_KEY_REQ"
	case GW_CS_GENERATE_NEW_KEY_CFM:
		return "GW_CS_GENERATE_NEW_KEY_CFM"
	case GW_CS_GENERATE_NEW_KEY_NTF:
		return "GW_CS_GENERATE_NEW_KEY_NTF"
	case GW_CS_REPAIR_KEY_REQ:
		return "GW_CS_REPAIR_KEY_REQ"
	case GW_CS_REPAIR_KEY_CFM:
		return "GW_CS_REPAIR_KEY_CFM"
	case GW_CS_REPAIR_KEY_NTF:
		return "GW_CS_REPAIR_KEY_NTF"
	case GW_CS_ACTIVATE_CONFIGURATION_MODE_REQ:
		return "GW_CS_ACTIVATE_CONFIGURATION_MODE_REQ"
	case GW_CS_ACTIVATE_CONFIGURATION_MODE_CFM:
		return "GW_CS_ACTIVATE_CONFIGURATION_MODE_CFM"
	case GW_GET_NODE_INFORMATION_REQ:
		return "GW_GET_NODE_INFORMATION_REQ"
	case GW_GET_NODE_INFORMATION_CFM:
		return "GW_GET_NODE_INFORMATION_CFM"
	case GW_GET_NODE_INFORMATION_NTF:
		return "GW_GET_NODE_INFORMATION_NTF"
	case GW_GET_ALL_NODES_INFORMATION_REQ:
		return "GW_GET_ALL_NODES_INFORMATION_REQ"
	case GW_GET_ALL_NODES_INFORMATION_CFM:
		return "GW_GET_ALL_NODES_INFORMATION_CFM"
	case GW_GET_ALL_NODES_INFORMATION_NTF:
		return "GW_GET_ALL_NODES_INFORMATION_NTF"
	case GW_GET_ALL_NODES_INFORMATION_FINISHED_NTF:
		return "GW_GET_ALL_NODES_INFORMATION_FINISHED_NTF"
	case GW_SET_NODE_VARIATION_REQ:
		return "GW_SET_NODE_VARIATION_REQ"
	case GW_SET_NODE_VARIATION_CFM:
		return "GW_SET_NODE_VARIATION_CFM"
	case GW_SET_NODE_NAME_REQ:
		return "GW_SET_NODE_NAME_REQ"
	case GW_SET_NODE_NAME_CFM:
		return "GW_SET_NODE_NAME_CFM"
	case GW_NODE_INFORMATION_CHANGED_NTF:
		return "GW_NODE_INFORMATION_CHANGED_NTF"
	case GW_NODE_STATE_POSITION_CHANGED_NTF:
		return "GW_NODE_STATE_POSITION_CHANGED_NTF"
	case GW_SET_NODE_ORDER_AND_PLACEMENT_REQ:
		return "GW_SET_NODE_ORDER_AND_PLACEMENT_REQ"
	case GW_SET_NODE_ORDER_AND_PLACEMENT_CFM:
		return "GW_SET_NODE_ORDER_AND_PLACEMENT_CFM"
	case GW_GET_GROUP_INFORMATION_REQ:
		return "GW_GET_GROUP_INFORMATION_REQ"
	case GW_GET_GROUP_INFORMATION_CFM:
		return "GW_GET_GROUP_INFORMATION_CFM"
	case GW_GET_GROUP_INFORMATION_NTF:
		return "GW_GET_GROUP_INFORMATION_NTF"
	case GW_SET_GROUP_INFORMATION_REQ:
		return "GW_SET_GROUP_INFORMATION_REQ"
	case GW_SET_GROUP_INFORMATION_CFM:
		return "GW_SET_GROUP_INFORMATION_CFM"
	case GW_GROUP_INFORMATION_CHANGED_NTF:
		return "GW_GROUP_INFORMATION_CHANGED_NTF"
	case GW_DELETE_GROUP_REQ:
		return "GW_DELETE_GROUP_REQ"
	case GW_DELETE_GROUP_CFM:
		return "GW_DELETE_GROUP_CFM"
	case GW_NEW_GROUP_REQ:
		return "GW_NEW_GROUP_REQ"
	case GW_NEW_GROUP_CFM:
		return "GW_NEW_GROUP_CFM"
	case GW_GET_ALL_GROUPS_INFORMATION_REQ:
		return "GW_GET_ALL_GROUPS_INFORMATION_REQ"
	case GW_GET_ALL_GROUPS_INFORMATION_CFM:
		return "GW_GET_ALL_GROUPS_INFORMATION_CFM"
	case GW_GET_ALL_GROUPS_INFORMATION_NTF:
		return "GW_GET_ALL_GROUPS_INFORMATION_NTF"
	case GW_GET_ALL_GROUPS_INFORMATION_FINISHED_NTF:
		return "GW_GET_ALL_GROUPS_INFORMATION_FINISHED_NTF"
	case GW_GROUP_DELETED_NTF:
		return "GW_GROUP_DELETED_NTF"
	case GW_HOUSE_STATUS_MONITOR_ENABLE_REQ:
		return "GW_HOUSE_STATUS_MONITOR_ENABLE_REQ"
	case GW_HOUSE_STATUS_MONITOR_ENABLE_CFM:
		return "GW_HOUSE_STATUS_MONITOR_ENABLE_CFM"
	case GW_HOUSE_STATUS_MONITOR_DISABLE_REQ:
		return "GW_HOUSE_STATUS_MONITOR_DISABLE_REQ"
	case GW_HOUSE_STATUS_MONITOR_DISABLE_CFM:
		return "GW_HOUSE_STATUS_MONITOR_DISABLE_CFM"
	case GW_COMMAND_SEND_REQ:
		return "GW_COMMAND_SEND_REQ"
	case GW_COMMAND_SEND_CFM:
		return "GW_COMMAND_SEND_CFM"
	case GW_COMMAND_RUN_STATUS_NTF:
		return "GW_COMMAND_RUN_STATUS_NTF"
	case GW_COMMAND_REMAINING_TIME_NTF:
		return "GW_COMMAND_REMAINING_TIME_NTF"
	case GW_SESSION_FINISHED_NTF:
		return "GW_SESSION_FINISHED_NTF"
	case GW_STATUS_REQUEST_REQ:
		return "GW_STATUS_REQUEST_REQ"
	case GW_STATUS_REQUEST_CFM:
		return "GW_STATUS_REQUEST_CFM"
	case GW_STATUS_REQUEST_NTF:
		return "GW_STATUS_REQUEST_NTF"
	case GW_WINK_SEND_REQ:
		return "GW_WINK_SEND_REQ"
	case GW_WINK_SEND_CFM:
		return "GW_WINK_SEND_CFM"
	case GW_WINK_SEND_NTF:
		return "GW_WINK_SEND_NTF"
	case GW_SET_LIMITATION_REQ:
		return "GW_SET_LIMITATION_REQ"
	case GW_SET_LIMITATION_CFM:
		return "GW_SET_LIMITATION_CFM"
	case GW_GET_LIMITATION_STATUS_REQ:
		return "GW_GET_LIMITATION_STATUS_REQ"
	case GW_GET_LIMITATION_STATUS_CFM:
		return "GW_GET_LIMITATION_STATUS_CFM"
	case GW_LIMITATION_STATUS_NTF:
		return "GW_LIMITATION_STATUS_NTF"
	case GW_MODE_SEND_REQ:
		return "GW_MODE_SEND_REQ"
	case GW_MODE_SEND_CFM:
		return "GW_MODE_SEND_CFM"
	case GW_MODE_SEND_NTF:
		return "GW_MODE_SEND_NTF"
	case GW_INITIALIZE_SCENE_REQ:
		return "GW_INITIALIZE_SCENE_REQ"
	case GW_INITIALIZE_SCENE_CFM:
		return "GW_INITIALIZE_SCENE_CFM"
	case GW_INITIALIZE_SCENE_NTF:
		return "GW_INITIALIZE_SCENE_NTF"
	case GW_INITIALIZE_SCENE_CANCEL_REQ:
		return "GW_INITIALIZE_SCENE_CANCEL_REQ"
	case GW_INITIALIZE_SCENE_CANCEL_CFM:
		return "GW_INITIALIZE_SCENE_CANCEL_CFM"
	case GW_RECORD_SCENE_REQ:
		return "GW_RECORD_SCENE_REQ"
	case GW_RECORD_SCENE_CFM:
		return "GW_RECORD_SCENE_CFM"
	case GW_RECORD_SCENE_NTF:
		return "GW_RECORD_SCENE_NTF"
	case GW_DELETE_SCENE_REQ:
		return "GW_DELETE_SCENE_REQ"
	case GW_DELETE_SCENE_CFM:
		return "GW_DELETE_SCENE_CFM"
	case GW_RENAME_SCENE_REQ:
		return "GW_RENAME_SCENE_REQ"
	case GW_RENAME_SCENE_CFM:
		return "GW_RENAME_SCENE_CFM"
	case GW_GET_SCENE_LIST_REQ:
		return "GW_GET_SCENE_LIST_REQ"
	case GW_GET_SCENE_LIST_CFM:
		return "GW_GET_SCENE_LIST_CFM"
	case GW_GET_SCENE_LIST_NTF:
		return "GW_GET_SCENE_LIST_NTF"
	case GW_GET_SCENE_INFOAMATION_REQ:
		return "GW_GET_SCENE_INFOAMATION_REQ"
	case GW_GET_SCENE_INFOAMATION_CFM:
		return "GW_GET_SCENE_INFOAMATION_CFM"
	case GW_GET_SCENE_INFOAMATION_NTF:
		return "GW_GET_SCENE_INFOAMATION_NTF"
	case GW_ACTIVATE_SCENE_REQ:
		return "GW_ACTIVATE_SCENE_REQ"
	case GW_ACTIVATE_SCENE_CFM:
		return "GW_ACTIVATE_SCENE_CFM"
	case GW_STOP_SCENE_REQ:
		return "GW_STOP_SCENE_REQ"
	case GW_STOP_SCENE_CFM:
		return "GW_STOP_SCENE_CFM"
	case GW_SCENE_INFORMATION_CHANGED_NTF:
		return "GW_SCENE_INFORMATION_CHANGED_NTF"
	case GW_ACTIVATE_PRODUCTGROUP_REQ:
		return "GW_ACTIVATE_PRODUCTGROUP_REQ"
	case GW_ACTIVATE_PRODUCTGROUP_CFM:
		return "GW_ACTIVATE_PRODUCTGROUP_CFM"
	case GW_ACTIVATE_PRODUCTGROUP_NTF:
		return "GW_ACTIVATE_PRODUCTGROUP_NTF"
	case GW_GET_CONTACT_INPUT_LINK_LIST_REQ:
		return "GW_GET_CONTACT_INPUT_LINK_LIST_REQ"
	case GW_GET_CONTACT_INPUT_LINK_LIST_CFM:
		return "GW_GET_CONTACT_INPUT_LINK_LIST_CFM"
	case GW_SET_CONTACT_INPUT_LINK_REQ:
		return "GW_SET_CONTACT_INPUT_LINK_REQ"
	case GW_SET_CONTACT_INPUT_LINK_CFM:
		return "GW_SET_CONTACT_INPUT_LINK_CFM"
	case GW_REMOVE_CONTACT_INPUT_LINK_REQ:
		return "GW_REMOVE_CONTACT_INPUT_LINK_REQ"
	case GW_REMOVE_CONTACT_INPUT_LINK_CFM:
		return "GW_REMOVE_CONTACT_INPUT_LINK_CFM"
	case GW_GET_ACTIVATION_LOG_HEADER_REQ:
		return "GW_GET_ACTIVATION_LOG_HEADER_REQ"
	case GW_GET_ACTIVATION_LOG_HEADER_CFM:
		return "GW_GET_ACTIVATION_LOG_HEADER_CFM"
	case GW_CLEAR_ACTIVATION_LOG_REQ:
		return "GW_CLEAR_ACTIVATION_LOG_REQ"
	case GW_CLEAR_ACTIVATION_LOG_CFM:
		return "GW_CLEAR_ACTIVATION_LOG_CFM"
	case GW_GET_ACTIVATION_LOG_LINE_REQ:
		return "GW_GET_ACTIVATION_LOG_LINE_REQ"
	case GW_GET_ACTIVATION_LOG_LINE_CFM:
		return "GW_GET_ACTIVATION_LOG_LINE_CFM"
	case GW_ACTIVATION_LOG_UPDATED_NTF:
		return "GW_ACTIVATION_LOG_UPDATED_NTF"
	case GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_REQ:
		return "GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_REQ"
	case GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_NTF:
		return "GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_NTF"
	case GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_CFM:
		return "GW_GET_MULTIPLE_ACTIVATION_LOG_LINES_CFM"
	case GW_SET_UTC_REQ:
		return "GW_SET_UTC_REQ"
	case GW_SET_UTC_CFM:
		return "GW_SET_UTC_CFM"
	case GW_RTC_SET_TIME_ZONE_REQ:
		return "GW_RTC_SET_TIME_ZONE_REQ"
	case GW_RTC_SET_TIME_ZONE_CFM:
		return "GW_RTC_SET_TIME_ZONE_CFM"
	case GW_GET_LOCAL_TIME_REQ:
		return "GW_GET_LOCAL_TIME_REQ"
	case GW_GET_LOCAL_TIME_CFM:
		return "GW_GET_LOCAL_TIME_CFM"
	case GW_PASSWORD_ENTER_REQ:
		return "GW_PASSWORD_ENTER_REQ"
	case GW_PASSWORD_ENTER_CFM:
		return "GW_PASSWORD_ENTER_CFM"
	case GW_PASSWORD_CHANGE_REQ:
		return "GW_PASSWORD_CHANGE_REQ"
	case GW_PASSWORD_CHANGE_CFM:
		return "GW_PASSWORD_CHANGE_CFM"
	case GW_PASSWORD_CHANGE_NTF:
		return "GW_PASSWORD_CHANGE_NTF"
	default:
		return fmt.Sprintf("<%d>", cmd)
	}
}
