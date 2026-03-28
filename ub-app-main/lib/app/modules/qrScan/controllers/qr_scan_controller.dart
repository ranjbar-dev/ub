import 'package:ai_barcode/ai_barcode.dart';
import 'package:file_picker/file_picker.dart';
import 'package:get/get.dart';
import 'package:qr_code_tools/qr_code_tools.dart';

import '../../../../utils/commonUtils.dart';
import '../../../../utils/logger.dart';

class QrScanController extends GetxController {
  final log = UBLogger.log;
  ScannerController scanner;
  final scannedValue = ''.obs;
  final deniedCameraPermission = false.obs;

  @override
  void onInit() {
    scanner = ScannerController(scannerResult: (v) {
      handleScanned(v);
    });
    super.onInit();
  }

  @override
  void onReady() async {
    super.onReady();
  }

  @override
  void onClose() {
    scanner.stopCamera();
    scanner.stopCameraPreview();
  }

  void startCameraUsage() async {
    await checkCameraPermission(
        onGranted: () => _startScanning(),
        onDenied: () {
          deniedCameraPermission.value = true;
        });
  }

  void handleScanned(String v) {
    scannedValue.value = v.replaceAll(' ', '');
    scanner.stopCamera();
    scanner.stopCameraPreview();
    if (v != null) {
      if (v.toLowerCase().contains('permission to access the camera')) {
        log.e('camera Permission Denied');
        deniedCameraPermission.value = true;
        return;
      }
    }
    Get.back(result: v);
  }

  void _startScanning() {
    Future.delayed(100.milliseconds).then((value) {
      scanner.startCamera();
      Future.delayed(100.milliseconds).then((value) {
        scanner.startCameraPreview();
      });
    });
  }

  void handleBrowsClick() async {
    await checkGalleryPermission(
        onGranted: () async {
          FilePickerResult result = await FilePicker.platform.pickFiles(
            allowMultiple: false,
            type: FileType.custom,
            allowedExtensions: ['jpg', 'png', 'jpeg'],
          );
          if (result != null) {
            final file = result.files[0].path;
            try {
              String data = await QrCodeToolsPlugin.decodeFrom(file);
              Get.back(result: data);
            } catch (e) {
              print(e.toString());
            }
          }
        },
        onDenied: () {});
  }
}
