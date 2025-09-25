#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/jiffies.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Alex Orr");
MODULE_DESCRIPTION("Modulo kernel para exponer PID, nombre y CPU en JSON");
MODULE_VERSION("1.0");

#define PROC_NAME "proc_monitor"

// Función que genera el contenido JSON para /proc
static int proc_show(struct seq_file *m, void *v) {
    struct task_struct *task;
    int first = 1;

    seq_puts(m, "[\n");
    for_each_process(task) {
        unsigned long long cpu_sec = (unsigned long long)(task->utime + task->stime) / HZ;

        if (!first)
            seq_puts(m, ",\n");
        first = 0;

        seq_printf(m,
            "  { \"pid\": %d, \"name\": \"%s\", \"cmdline\": \"N/A\", \"cpu_sec\": %llu, \"rss\": 0 }",
            task->pid,
            task->comm,
            cpu_sec
        );
    }
    seq_puts(m, "\n]\n");
    return 0;
}

// Abrir /proc
static int proc_open(struct inode *inode, struct file *file) {
    return single_open(file, proc_show, NULL);
}

// Operaciones de /proc
static const struct proc_ops proc_fops = {
    .proc_open    = proc_open,
    .proc_read    = seq_read,
    .proc_lseek   = seq_lseek,
    .proc_release = single_release,
};

// Inicialización del módulo
static int __init proc_init(void) {
    proc_create(PROC_NAME, 0, NULL, &proc_fops);
    printk(KERN_INFO "proc_monitor cargado\n");
    return 0;
}

// Limpieza del módulo
static void __exit proc_exit(void) {
    remove_proc_entry(PROC_NAME, NULL);
    printk(KERN_INFO "proc_monitor descargado\n");
}

module_init(proc_init);
module_exit(proc_exit);
