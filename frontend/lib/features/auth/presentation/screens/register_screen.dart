import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../../../core/widgets/responsive_layout.dart';
import '../providers/auth_provider.dart';
import 'package:go_router/go_router.dart';

class RegisterScreen extends StatelessWidget {
  const RegisterScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ResponsiveLayout(
        mobileBody: const _RegisterForm(),
        webBody: Center(child: SizedBox(width: 400, child: _RegisterForm())),
      ),
    );
  }
}

class _RegisterForm extends StatefulWidget {
  const _RegisterForm();

  @override
  State<_RegisterForm> createState() => _RegisterFormState();
}

class _RegisterFormState extends State<_RegisterForm> {
  final _formKey = GlobalKey<FormState>();
  String firstName = '';
  String lastName = '';
  String email = '';
  String password = '';
  String username = '';
  bool isLoading = false;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final size = MediaQuery.of(context).size;
    final buttonFontSize = 16.0;
    final borderRadius = BorderRadius.circular(10);

    return Padding(
      padding: EdgeInsets.symmetric(
        horizontal: size.width < 800 ? 24 : 0,
        vertical: 32,
      ),
      child: Form(
        key: _formKey,
        child: ListView(
          shrinkWrap: true,
          children: [
            Text(
              "Create An Account",
              style: TextStyle(
                fontFamily: 'Montserrat',
                fontWeight: FontWeight.bold,
                fontSize: 24,
                color: colorScheme.primary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 24),
            TextFormField(
              decoration: InputDecoration(
                labelText: "First Name",
                border: OutlineInputBorder(borderRadius: borderRadius),
              ),
              onChanged: (v) => firstName = v,
              validator: (v) => v != null && v.isNotEmpty ? null : "Required",
            ),
            const SizedBox(height: 12),
            TextFormField(
              decoration: InputDecoration(
                labelText: "Last Name",
                border: OutlineInputBorder(borderRadius: borderRadius),
              ),
              onChanged: (v) => lastName = v,
              validator: (v) => v != null && v.isNotEmpty ? null : "Required",
            ),
            const SizedBox(height: 12),
            TextFormField(
              decoration: InputDecoration(
                labelText: "Username",
                border: OutlineInputBorder(borderRadius: borderRadius),
              ),
              onChanged: (v) => username = v,
              validator: (v) => v != null && v.isNotEmpty ? null : "Required",
            ),
            const SizedBox(height: 12),
            TextFormField(
              decoration: InputDecoration(
                labelText: "Email",
                border: OutlineInputBorder(borderRadius: borderRadius),
              ),
              onChanged: (v) => email = v,
              keyboardType: TextInputType.emailAddress,
              validator: (v) =>
                  v != null && v.contains('@') ? null : "Enter a valid email",
            ),
            const SizedBox(height: 12),
            TextFormField(
              decoration: InputDecoration(
                labelText: "Password",
                border: OutlineInputBorder(borderRadius: borderRadius),
              ),
              obscureText: true,
              onChanged: (v) => password = v,
              validator: (v) =>
                  v != null && v.length >= 6 ? null : "Min 6 chars",
            ),
            const SizedBox(height: 20),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: colorScheme.primary,
                  foregroundColor: colorScheme.onPrimary,
                  textStyle: TextStyle(
                    fontFamily: 'Montserrat',
                    fontWeight: FontWeight.bold,
                    fontSize: buttonFontSize,
                  ),
                  shape: RoundedRectangleBorder(borderRadius: borderRadius),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                onPressed: isLoading
                    ? null
                    : () async {
                        if (_formKey.currentState?.validate() ?? false) {
                          setState(() => isLoading = true);
                          await Provider.of<AuthProvider>(
                            context,
                            listen: false,
                          ).register(
                            name: lastName,
                            firstname: firstName,
                            username: username,
                            email: email,
                            password: password,
                            context: context,
                          );
                          setState(() => isLoading = false);
                        }
                      },
                child: isLoading
                    ? const CircularProgressIndicator()
                    : const Text("Sign Up"),
              ),
            ),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                icon: const Icon(Icons.g_mobiledata, color: Colors.blue),
                label: const Text("Continue with Google"),
                style: OutlinedButton.styleFrom(
                  textStyle: TextStyle(
                    fontFamily: 'Montserrat',
                    fontWeight: FontWeight.bold,
                    fontSize: buttonFontSize,
                  ),
                  side: BorderSide(color: colorScheme.primary),
                  shape: RoundedRectangleBorder(borderRadius: borderRadius),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                onPressed: () {
                  // TODO: Google sign up
                },
              ),
            ),
            const SizedBox(height: 12),
            Text(
              "By signing up, you agree with the Terms of Service and Privacy Policy",
              style: TextStyle(
                fontFamily: 'Montserrat',
                fontSize: 12,
                color: colorScheme.secondary,
              ),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: OutlinedButton(
                style: OutlinedButton.styleFrom(
                  textStyle: TextStyle(
                    fontFamily: 'Montserrat',
                    fontWeight: FontWeight.bold,
                    fontSize: buttonFontSize,
                  ),
                  shape: RoundedRectangleBorder(borderRadius: borderRadius),
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                onPressed: () {
                  context.go('/login');
                },
                child: const Text("Already have an account?"),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
